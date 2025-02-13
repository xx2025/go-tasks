package executor

import (
	"context"
	"errors"
	"fmt"
	"go-tasks/grpc/service"
	"go-tasks/internal/rpc"
	"go-tasks/utils"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

var ProcessMap sync.Map

func LoadOrSetProcessMap(processId int64, process *Process) (*Process, bool) {
	value, ok := ProcessMap.LoadOrStore(processId, process)

	if ok {
		return value.(*Process), true
	}
	return process, false
}

func SetProcessMap(processId int64, process *Process) {
	ProcessMap.Store(processId, process)
}

func GetProcessMap(processId int64) (*Process, error) {

	if value, ok := ProcessMap.Load(processId); ok {
		return value.(*Process), nil
	}

	return nil, errors.New("process not exist")
}

func DelProcessMap(processId int64) {
	ProcessMap.Delete(processId)
}

type Process struct {
	cmd        *exec.Cmd
	processId  int64
	command    string
	args       []string
	maxRetries int
	mu         sync.Mutex
	stopChan   chan struct{} //主动停止进程信号
	doneChan   chan struct{} //进程退出，自动重启信号
}

func NewProcess(processId int64, command string, args []string, maxRetries int) (*Process, error) {

	p := &Process{
		//cmd:      exec.Command(command, args...),
		processId:  processId,
		command:    command,
		args:       args,
		maxRetries: maxRetries,
		stopChan:   make(chan struct{}),
		doneChan:   make(chan struct{}),
	}

	_, exists := LoadOrSetProcessMap(processId, p)
	if exists {
		return nil, errors.New("该进程已创建")
	}

	return p, nil
}

func (p *Process) createCommand() {
	p.cmd = exec.Command(p.command, p.args...)
}

func (p *Process) Start() {
	defer DelProcessMap(p.processId)
	for {

		p.mu.Lock()

		success := false
		for i := 0; i <= p.maxRetries; i++ {
			p.createCommand()
			err := p.cmd.Start()
			if err == nil {
				success = true
				break
			}
		}

		p.mu.Unlock()
		if !success {
			//启动失败
			//通知调度中心， 进程启动失败
			p.processNotice("进程启动失败")
			return
		}
		SetProcessMap(p.processId, p)

		go p.waitForProcess()

		//通知调度中心， 进程启动了
		p.processNotice("进程已启动")
		select {
		case <-p.doneChan:
			//进程异常退出信号, 重新启动

		case <-p.stopChan:
			//收到进程停止信号
			p.stopProcess()
			return
		}
	}
}

// 接收外部信号， 停止进程
func (p *Process) Stop() {
	if p.cmd.Process == nil {
		return
	}
	p.mu.Lock()
	p.stopChan <- struct{}{}
	p.mu.Unlock()
}

func (p *Process) waitForProcess() {
	err := p.cmd.Wait()
	if err != nil {
		fmt.Printf("Process exited with error: %v\n", err)
	}
	p.mu.Lock()
	p.doneChan <- struct{}{} //进程退出通知
	p.mu.Unlock()
}

func (p *Process) stopProcess() {
	if p.cmd.Process == nil {
		return
	}
	p.mu.Lock()
	defer func() {
		close(p.stopChan)
		p.mu.Unlock()
	}()
	err := p.cmd.Process.Signal(syscall.SIGTERM)
	if err != nil {
		fmt.Printf("Failed to stop process: %v\n", err)
	}
	// 等待一段时间，检查进程是否退出
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	go func() {
		<-ctx.Done()
		// 超时后发送 SIGKILL
		_ = p.cmd.Process.Kill()
	}()
	// 等待进程退出
	_ = p.cmd.Wait()
	p.processNotice("进程已退出")
}

func (p *Process) processNotice(msg string) {
	f := func(ctx context.Context, client service.MasterClient) (interface{}, error) {
		return client.ProcessNotice(ctx, &service.ProcessNoticeRequest{
			ProcessId: p.processId,
			Message:   msg,
		})
	}
	_, err := rpc.ExecuteGRPCRequest(f)

	if err != nil {
		utils.Logger().Error("processNoticeError: " + err.Error())
	}
}
