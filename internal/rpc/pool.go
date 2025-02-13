package rpc

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"sync"
	"time"
)

type Client struct {
	Conn       *grpc.ClientConn
	createTime time.Time
	putTime    time.Time
}

type GRPCPool struct {
	address          string
	maxIdleConnCount int           //最大空闲数
	maxLifetime      time.Duration //最大生存时长
	maxIdleTime      time.Duration //最大空闲时长
	pool             chan *Client
	mu               sync.Mutex
}

func newGRPCPool(address string) *GRPCPool {
	maxIdleConnCount := 10
	pool := make(chan *Client, maxIdleConnCount*2)
	p := &GRPCPool{
		address:          address,
		maxIdleConnCount: 10,
		maxLifetime:      time.Second * 5 * 3600,
		maxIdleTime:      time.Second * 1 * 3600,
		pool:             pool,
	}
	return p
}

func (p *GRPCPool) initPoolConn(num int) {
	for i := 1; i <= num; i++ {
		c, err := p.newConn()
		if err != nil {
			continue
		}

		p.put(c)

	}
}

func (p *GRPCPool) get() (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	for {
		select {
		case c := <-p.pool:
			if c != nil {
				valid := c.checkConnValid(p.maxLifetime, p.maxIdleTime)
				if !valid {
					c.close()
				} else {
					return c, nil
				}
			}
		case <-ctx.Done():
			c, err := p.newConn()
			if err != nil {
				return nil, errors.New("获取连接失败")
			}
			return c, nil
		}
	}
}

func (p *GRPCPool) put(c *Client) {
	if len(p.pool) >= p.maxIdleConnCount {
		c.close()
		return
	}
	if !c.checkConnValid(p.maxLifetime, p.maxIdleTime) {
		c.close()
		return
	}
	c.putTime = time.Now()
	p.pool <- c
}

func (p *GRPCPool) newConn() (*Client, error) {
	kaParams := keepalive.ClientParameters{
		Time:                60 * time.Second,
		Timeout:             time.Second,
		PermitWithoutStream: true,
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	conn, err := grpc.NewClient(
		p.address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(kaParams),
	)
	if err != nil {

		return nil, err
	}

	return &Client{
		Conn:       conn,
		createTime: time.Now(),
	}, nil
}

func (c *Client) checkConnValid(maxLifetime, maxIdleTime time.Duration) bool {
	if c.Conn.GetState() == connectivity.Shutdown || c.Conn.GetState() == connectivity.TransientFailure {
		return false
	}
	if c.createTime.Add(maxLifetime).Before(time.Now()) {
		return false
	}
	if c.putTime.Add(maxIdleTime).Before(time.Now()) {
		return false
	}
	return true
}

func (c *Client) close() {
	err := c.Conn.Close()
	if err != nil {
		return
	}
}
