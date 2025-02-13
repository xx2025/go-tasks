## go-tasks
基于golang开发的分布式任务调度中心、进程管理中心

## 项目特色

- **简洁易用**：多项目、多节点部署，基于golang二进制安装，无需配置运行环境

- **任务统一管理，多节点执行**：所有任务通过web界面统一管理，然后调度到各执行器执行，免去cron等类似服务的复杂操作。

- **内置进程管理**：内置守护进程管理功能，轻松应对若干异步进程需求。

![架构模式](./resource/system/go-tasks.png)

## WEB界面图
<img src="./resource/system/node.png" width="400px" alt="图片描述"> 
<img src="./resource/system/project.png" width="400px"> 
<img src="./resource/system/task.png" width="200px"> 
<img src="./resource/system/task_logs.png" width="200px">
<img src="./resource/system/process.png" width="200px">
<img src="./resource/system/process_detail.png" width="200px">

## WEB界面图
![node](./resource/system/node.png)![project](./resource/system/project.png)
![task](./resource/system/task.png)
![task-logs](./resource/system/task_logs.png)
![process](./resource/system/process.png)
![process_detail](./resource/system/process_detail.png)

