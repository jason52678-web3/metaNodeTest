# blog-system
博客系统后端（gin & gorm）

# 概要
这是一个基于 Gin + GORM 的简单博客系统后端。包含用户注册、登录（JWT）、文章和评论的增删改查接口。

运行环境
Go 1.25.4
MySQL

# 仓库结构（关键目录）
- `conf/` - 系统启动时需要的配置文件
- `controller/` -  服务的入口，负责处理路由、参数校验，请求转发
- `dao/` -  处理数据和存储相关功能
- `logger/` -  系统日志处理
- `logic/`  - 逻辑层，负责处理业务逻辑
- `middlewares/`  - 中间件处理（token存储）
- `models/`  - 相关参数的数据结构
- `pkg/`  - 第三方库
- `router/`  - 路由处理
- `settings/` -  配置文件参数动态加载
- `main.go`   -  系统入口，总体调用流程
- `.air.toml`  - air配置文件

# 配置（config.yaml）
## 在conf目录下包含 config.yaml，示例内容：
`name: "MyBlog"`

`mode: "dev"`

`port: 8088`

`start_time: "2025-11-26"`

`machine_id: 1`

# log:
level: "debug"

filename : "MyBlog.log"

max_size: 200

max_age: 30

max_backups: 7

# mysql:
host: "127.0.0.1"

port: 3306

user: "root"

password: "123456"

dbname: "myblogdb"

max_open_conns: 20

max_idle_conns: 10

# redis:
host: "127.0.0.1"

port: 16379

password: ""

db: 0

pool_size: 10

# auth:
jwt_expire: 8760

# 依赖安装 & 编译
`请根据本地环境修改 mysql 配置（用户名、密码、数据库名等）。`

`在项目根目录运行：`

`go mod download`
`go build -o blog-system ./`

`也可以直接用 go run 运行（适合开发）：`

`go run main.go`

注意：程序会在启动时根据 config.yaml 初始化数据库连接
也支持外部指定的合法yamf配置文件，比如：blog-system new_config.yaml

# HTTP 接口 & 测试（curl）
路由在 `router/router.go` 中定义，主要接口如下：

1) 公开接口（无需 JWT）
- POST /api/v1/signup — 用户注册
- POST /api/v1/login — 用户登录（返回 token）
- GET  /api/v1/postsall - 浏览所有帖子列表
- GET  /api/v1/postsdetail/:title -浏览某篇文章的详情

2) 受保护接口（需携带 token，支持 Authorization: Bearer）
- POST /api/v1/createpost — 新建文章
- POST /api/v1/updatepost — 更新文章
- POST /api/v1/deletepost/ — 删除文章
- POST /api/v1/createcomment — 新增评论
- GET  /api/v1/getpostcomments/:post_id" — 查询评论（按 postID）

## 如何运行并收集真实测试结果
1) 确保MySQL可用，并在 config.yaml 中配置正确的 username/password/db-name。
2) 如果数据库中尚无myblogdb数据库，先创建：
- `CREATE DATABASE myblogdb CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;`

3) 启动服务：
```bash
go run main.go
```
4) 开启Postman客户端进行GET/POST测试
- Postman客户端下载地址：https://www.postman.com/downloads/

注意事项与常见问题
- 如果遇到token相关错误，请确认登录确实返回了 token 并在请求中以 Authorization: Bearer <token> 形式传入。
- 如果启动时报 DB 连接错误，检查 config.yaml 的 MySQL 字段并确保 MySQL 可用。