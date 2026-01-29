# 衣搭库后端服务

衣搭库（YiDaKu）是一款基于 AI 视觉识别的智能穿搭助手后端服务。

## 技术栈

- **语言**: Go 1.21+
- **框架**: Gin
- **数据库**: MySQL 8.0+
- **缓存**: Redis（可选）
- **存储**: 七牛云 OSS
- **认证**: JWT

## 项目结构

```
server/
├── cmd/
│   └── server/
│       └── main.go           # 主入口
├── internal/
│   ├── config/               # 配置管理
│   ├── handler/              # HTTP处理器
│   ├── middleware/           # 中间件
│   ├── model/                # 数据模型
│   ├── repository/           # 数据访问层
│   ├── service/              # 业务逻辑层
│   └── adapter/              # 外部服务适配器
├── pkg/
│   ├── auth/                 # JWT认证
│   ├── cache/                # Redis缓存
│   ├── response/             # 统一响应格式
│   └── storage/              # 七牛云存储
├── migrations/               # 数据库迁移脚本
├── config.example.yaml       # 配置文件示例
├── go.mod
└── README.md
```

## 快速开始

### 1. 环境准备

- Go 1.21+
- MySQL 8.0+
- Redis（可选，用于缓存）

### 2. 配置

复制配置文件并修改：

```bash
cp config.example.yaml config.yaml
```

编辑 `config.yaml` 配置数据库、Redis、七牛云等信息。

### 3. 初始化数据库

```bash
# 执行数据库迁移脚本
mysql -u root -p < migrations/001_init.sql
mysql -u root -p < migrations/002_seed_data.sql
```

### 4. 安装依赖

```bash
cd server
go mod tidy
```

### 5. 运行

```bash
go run cmd/server/main.go
# 或指定配置文件
go run cmd/server/main.go /path/to/config.yaml
```

服务将在配置的端口启动（默认 8080）。

### 6. 编译

```bash
# Windows
go build -o yidaiku-server.exe cmd/server/main.go

# Linux/Mac
go build -o yidaiku-server cmd/server/main.go
```

## API 接口

### 认证相关

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | /api/v1/auth/register | 用户注册 |
| POST | /api/v1/auth/login | 用户登录 |
| POST | /api/v1/auth/refresh | 刷新Token |

### 用户相关

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | /api/v1/user/profile | 获取用户信息 |
| PUT | /api/v1/user/profile | 更新用户信息 |
| PUT | /api/v1/user/password | 修改密码 |
| GET | /api/v1/user/preference | 获取偏好设置 |
| PUT | /api/v1/user/preference | 更新偏好设置 |

### 衣橱管理

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | /api/v1/wardrobe/items | 创建衣物 |
| GET | /api/v1/wardrobe/items | 获取衣物列表 |
| GET | /api/v1/wardrobe/items/:id | 获取衣物详情 |
| PUT | /api/v1/wardrobe/items/:id | 更新衣物 |
| DELETE | /api/v1/wardrobe/items/:id | 删除衣物 |
| PUT | /api/v1/wardrobe/items/:id/retired | 设置淘汰状态 |
| POST | /api/v1/wardrobe/items/:id/restore | 恢复已删除衣物 |
| PUT | /api/v1/wardrobe/items/batch/retired | 批量设置淘汰 |
| DELETE | /api/v1/wardrobe/items/batch | 批量删除 |
| GET | /api/v1/wardrobe/stats | 获取衣橱统计 |
| GET | /api/v1/wardrobe/categories | 获取品类列表 |
| GET | /api/v1/wardrobe/styles | 获取风格列表 |
| GET | /api/v1/wardrobe/seasons | 获取季节列表 |
| GET | /api/v1/wardrobe/colors | 获取颜色列表 |

### 穿搭推荐

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | /api/v1/outfits/recommend | 生成穿搭推荐 |
| POST | /api/v1/outfits/replace | 换一件 |
| GET | /api/v1/outfits/favorites | 获取收藏列表 |
| PUT | /api/v1/outfits/:id/favorite | 收藏/取消收藏 |
| POST | /api/v1/outfits/:id/apply | 标记为已应用 |
| GET | /api/v1/outfits/templates | 获取穿搭模板 |
| GET | /api/v1/outfits/occasions | 获取场合列表 |

### 文件上传

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | /api/v1/upload/init | 初始化分块上传 |
| POST | /api/v1/upload/chunk | 上传分块 |
| POST | /api/v1/upload/complete | 完成上传 |
| POST | /api/v1/upload/simple | 简单上传（小文件） |
| GET | /api/v1/upload/pending | 获取待处理任务 |
| GET | /api/v1/upload/progress/:task_id | 获取上传进度 |
| DELETE | /api/v1/upload/:task_id | 取消上传 |
| POST | /api/v1/upload/:task_id/retry | 重试上传 |

## 响应格式

所有接口返回统一的 JSON 格式：

```json
{
    "code": 0,
    "message": "success",
    "data": {}
}
```

分页数据格式：

```json
{
    "code": 0,
    "message": "success",
    "data": {
        "list": [],
        "total": 100,
        "page": 1,
        "page_size": 20
    }
}
```

## 错误码

| 错误码 | 描述 |
|--------|------|
| 0 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权/Token过期 |
| 403 | 禁止访问 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## 开发说明

### MVP 版本功能

当前 MVP 版本实现以下核心功能：

1. **用户系统**: 注册、登录、偏好设置
2. **衣橱管理**: 上传衣物、手动分类（品类/颜色/风格/季节）、筛选、淘汰管理
3. **穿搭推荐**: 基于规则引擎的推荐、换一件、收藏
4. **断点续传**: 支持大文件分块上传

### 后续版本规划

- AI 衣物识别（接入豆包 AI）
- AI 穿搭推荐增强
- 穿搭记录/日记
- 每日推送
- 数据导出

## License

MIT
