

为何IM业务用Go写很合适

---

### 分析

> 面向复杂度架构设计流程：
>
> 1. 分析 IM 业务特点 - 业务复杂度低，质量复杂度高，量化是连接数多，并发高
> 2. 分析 Go 的特点及本质 - Go 协程和并发模型适合 IM 业务

1. 业务分析

   - 业务复制度分析

     IM业务逻辑相对简单，质量复杂度高，但业务复杂度低

2. Go 本身分析

   - 快速开发
   - 协程和并发模型适合 IM 业务

### [开源IM系统 OpenIM](https://github.com/OpenIMSDK/Open-IM-Server)

![OpenIM 架构](https://tva1.sinaimg.cn/large/e6c9d24egy1h2q68agnzdj216x0u0dl2.jpg)