## 项目介绍
基于golang开发的分布式任务调度中心、进程管理中心

## 项目特色

- **简洁易用**：多项目、多节点部署，基于golang二进制安装，无需配置运行环境

- **任务统一管理，多节点执行**：所有任务通过web界面统一管理，然后调度到各执行器执行，免去cron等类似服务的复杂操作。

- **内置进程管理**：内置守护进程管理功能，轻松应对若干异步进程需求。

![架构模式](https://raw.githubusercontent.com/xx2025/go-tasks/refs/heads/main/resource/system/go-tasks.png)

## 启动方式
- **导入数据结构**
```
    go-tasks.sql， 账号密码：root:123456
    配置 master.config.toml
    配置 node.config.toml
```
    
- **启动master**
```
    cd ./cmd
    go run master.go
```
- **启动node**
```a
    cd ./cmd
    go run node.go
```
- **启动web**
```
    cd ./web
    pnpm install
    npm run dev
```

## WEB界面图
![node](https://raw.githubusercontent.com/xx2025/go-tasks/refs/heads/main/resource/system/node.png)
![project](https://raw.githubusercontent.com/xx2025/go-tasks/refs/heads/main/resource/system/project.png)
![task](https://raw.githubusercontent.com/xx2025/go-tasks/refs/heads/main/resource/system/task.png)
![task-logs](https://raw.githubusercontent.com/xx2025/go-tasks/refs/heads/main/resource/system/task_logs.png)
![process](https://raw.githubusercontent.com/xx2025/go-tasks/refs/heads/main/resource/system/process.png)
![process_detail](https://raw.githubusercontent.com/xx2025/go-tasks/refs/heads/main/resource/system/process_detail.png)

