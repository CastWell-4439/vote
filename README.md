项目简介
这是一个基于 Go 语言开发的投票系统后端服务，提供用户管理等基础功能，可扩展实现投票相关核心业务逻辑。目前已实现用户的注册、登录、信息查询、新增、更新、删除等功能

技术栈
编程语言：Go 1.25.1
Web 框架：Gin v1.11.0
ORM 工具：GORM v1.31.0
数据库：MySQL
缓存 / 会话存储：Redis
日志组件：标准库slog及自定义日志处理
其他依赖：go-sql-driver/mysql、gin-contrib/sessions


安装与部署
1. 克隆项目
   ```bash
   git clone <项目仓库地址>
   cd vote/vote
```

2. 配置环境
数据库配置
修改 config/database.go 中的数据库连接信息，适配你的 MySQL 环境：

```go
// config/database.go
username := "你的MySQL用户名"
password := "你的MySQL密码"
host := "MySQL主机地址"
port := 3306 
Dbname := "gorm" 
timeout := "10s"
```

Redis 配置
修改 config/redis.go 中的 Redis 地址

```go
// config/redis.go
const (
    RedisAddress = "localhost:6379" 
)
```

3. 安装依赖
```bash
go mod tidy
```

4. 启动服务
服务默认启动在 :8888 端口，可通过 http://localhost:8888 访问
```bash
go run start.go
```


目录结构

```plaintext
vote/
├── config/           # 配置文件（数据库、Redis等）
│   ├── database.go   # MySQL连接配置
│   └── redis.go      # Redis地址配置
├── controllers/      # 控制器（处理HTTP请求）
│   ├── Register.go   # 注册/登录相关接口
│   ├── user.go       # 用户CRUD接口
│   └── return.go     # 统一响应格式
├── dao/              # 数据访问层初始化
│   └── dao.go
├── logger/           # 日志处理
│   └── logger.go     # 日志写入、中间件等
├── model/            # 数据模型
│   └── user.go       # 用户模型及数据库操作
├── router/           # 路由配置
│   └── router.go     # 路由注册及中间件
├── start.go          # 程序入口
├── go.mod            
└── go.sum            
```

2025/10/14 更新
加入了shell脚本自动备份Mysql数据，源代码放在 /scripts/shell里了   

加入了lua脚本实现投票的控制，在 /config/redis.go文件中加入了连接池，并且创建了controller/controller.go文件用于和lua交互

lua源代码放在了 /scripts/lua中   

2025/10/16 更新
改了一点小错误，也修改了备份脚本

2025/10/16 再次更新
修改了logger和controller里面的致命错误（（   
特别感谢ddk愿意看我又臭又长的代码并帮我debug，你是我爹（


2025/10/19 更新  
1.设置了redis session的过期时间  
2.将注册逻辑中的明文密码改为了bcrypt  
3.统一了之前偷懒写的不一致的返回格式  
4.优化了backup脚本
5.上次以为没用上而删掉的实际上用到了，给加回来了


2025/10/20 更新  
1.重写了backup脚本，这次不再使用明文密码，而是使用配置文件和环境变量
2.意外发现了神秘的拼写错误，已修改（运行没问题oops








