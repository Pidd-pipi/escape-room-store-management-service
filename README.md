# Escape Room Ops 密室逃脱门店管理后端

密室逃脱门店运营平台的后端 API，提供主题房间、游戏场次、报名、逃脱记录、排行榜与营收分析等能力。

## 技术栈
- Go 1.22
- Gin + GORM
- PostgreSQL 15
- JWT + RBAC

## 标准命令
```bash
go build ./...        # 编译
go test ./...         # 运行测试
go run ./cmd/server   # 启动 HTTP 服务
```

## 环境变量
| 变量 | 说明 | 默认值 |
| --- | --- | --- |
| DB_HOST | PostgreSQL 主机 | localhost |
| DB_PORT | PostgreSQL 端口 | 5432 |
| DB_NAME | 数据库名 | escape_room_db |
| DB_USER | 数据库用户 | escape_room_user |
| DB_PASSWORD | 数据库密码 | escape_room_pwd |
| JWT_SECRET | JWT 签名密钥 | change_me_to_a_long_random_string |
| APP_CORS_ORIGINS | 允许跨域来源 | http://localhost:28504 |

## 目录结构
```
backend/
├── cmd/server/
└── internal/
    ├── config/       # 配置解析
    ├── model/        # 实体定义
    ├── repository/   # 数据访问层
    ├── service/      # 业务逻辑层
    ├── handler/      # HTTP 接口层
    ├── router/       # 路由注册
    ├── middleware/   # auth/rbac/rate_limiter/error_handler/request_id
    ├── dto/          # 请求/响应结构体
    ├── constants/    # 枚举、错误码、日志模板、文案
    └── util/         # jwt/logger/formatters/app_error/escape_rate_calculator
```
