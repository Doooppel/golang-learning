# 🚀 Golang Web Demo 项目

这是一个完整的Go语言Web开发学习项目，从基础语法到高级特性，涵盖了Web开发的核心概念。

## 📚 项目特点

- ✅ **完整的CRUD功能** - 用户增删改查
- ✅ **内存数据库** - 使用Map + 读写锁实现并发安全的内存存储
- ✅ **数据持久化** - 自动保存到JSON文件，重启后数据不丢失
- ✅ **RESTful API** - 标准的REST API设计
- ✅ **从基础到高级** - 包含大量Go语言语法示例
- ✅ **生产级代码** - 错误处理、参数验证、代码组织

## 🏗️ 项目结构

```
golang-web-demo/
├── main.go                    # 程序入口
├── go.mod                     # Go模块定义
├── data/                      # 数据文件目录（自动创建）
│   └── users.json            # 用户数据JSON文件（自动生成）
├── internal/                  # 内部包（不对外暴露）
│   ├── models/               # 数据模型
│   │   └── user.go           # 用户模型定义
│   ├── database/             # 数据访问层
│   │   └── memory_db.go      # 内存数据库实现（含持久化）
│   ├── handlers/             # HTTP处理器
│   │   └── user_handler.go   # 用户相关处理器
│   └── router/               # 路由配置
│       └── router.go         # 路由设置
├── examples/                  # 学习示例代码
│   ├── basic_syntax.go       # 基础语法示例
│   └── advanced_features.go  # 高级特性示例
└── README.md                 # 项目文档
```

## 🚀 快速开始

### 1. 安装依赖

```bash
go mod download
```

### 2. 运行项目

```bash
go run main.go
```

服务器将在 `http://localhost:8080` 启动

### 3. 访问API文档

打开浏览器访问：`http://localhost:8080/api/docs`

## 📖 API 端点

### 健康检查
```
GET /api/v1/health
```

### 用户管理

#### 创建用户
```bash
POST /api/v1/users
Content-Type: application/json

{
  "name": "张三",
  "email": "zhangsan@example.com",
  "age": 25
}
```

#### 获取所有用户
```bash
GET /api/v1/users
```

#### 获取单个用户
```bash
GET /api/v1/users/:id
```

#### 更新用户
```bash
PUT /api/v1/users/:id
Content-Type: application/json

{
  "name": "李四",
  "email": "lisi@example.com",
  "age": 30
}
```

#### 删除用户
```bash
DELETE /api/v1/users/:id
```

## 💾 数据持久化

项目使用JSON文件自动保存数据，数据存储在 `data/users.json` 文件中：

- **自动保存**：每次创建、更新或删除用户时，数据会自动保存到文件
- **自动加载**：服务器启动时，会自动从文件加载之前保存的数据
- **文件格式**：使用JSON格式，便于阅读和调试
- **原子写入**：使用临时文件+重命名的方式，确保数据完整性

### 数据文件位置
```
data/users.json
```

### 手动备份数据
如果需要备份数据，直接复制 `data/users.json` 文件即可。

### 重置数据
如果需要清空所有数据，删除 `data/users.json` 文件即可，服务器会在下次启动时创建新的空数据库。

## 🧪 测试示例

### 使用 curl 测试

```bash
# 创建用户
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Doppel","email":"zhangsan@example.com","age":25}'

# 获取所有用户
curl http://localhost:8080/api/v1/users

# 获取单个用户（假设ID为1）
curl http://localhost:8080/api/v1/users/1

# 更新用户
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"李四","age":30}'

# 删除用户
curl -X DELETE http://localhost:8080/api/v1/users/1
```

## 📝 学习要点

### 基础语法（examples/basic_syntax.go）

1. **变量声明** - var、短变量声明、类型推断
2. **数据类型** - 基本类型、数组、切片、Map
3. **结构体** - 定义和使用
4. **控制流** - if/else、switch、for、range
5. **函数** - 普通函数、多返回值、命名返回值
6. **错误处理** - Go语言的错误处理模式
7. **接口** - 接口定义和实现
8. **Goroutine** - 并发编程基础
9. **Channel** - 通道通信
10. **defer** - 延迟执行

### 高级特性（examples/advanced_features.go）

1. **接口和类型断言** - 多态实现
2. **泛型** - Go 1.18+ 泛型编程
3. **反射** - 运行时类型信息
4. **并发编程** - Goroutine、Channel、Select
5. **Context** - 上下文管理
6. **同步原语** - Mutex、RWMutex、WaitGroup
7. **自定义错误** - 错误类型定义
8. **函数式编程** - 函数作为一等公民

### Web开发实践

1. **项目组织** - internal包的使用
2. **依赖注入** - 通过构造函数注入依赖
3. **错误处理** - HTTP错误响应
4. **参数验证** - Gin的binding验证
5. **并发安全** - 使用sync.RWMutex保护共享数据
6. **RESTful设计** - 标准的REST API设计

## 🔧 技术栈

- **Go 1.21+** - 编程语言
- **Gin** - Web框架
- **标准库** - sync、time、fmt等

## 📚 学习路径建议

1. **第一步**：阅读 `examples/basic_syntax.go`，理解Go基础语法
2. **第二步**：查看 `internal/models/user.go`，了解结构体定义
3. **第三步**：研究 `internal/database/memory_db.go`，学习并发编程
4. **第四步**：阅读 `internal/handlers/user_handler.go`，理解HTTP处理
5. **第五步**：查看 `internal/router/router.go`，学习路由配置
6. **第六步**：阅读 `examples/advanced_features.go`，掌握高级特性

## 🧪 测试

项目包含完整的测试套件，覆盖所有功能：

```bash
# 运行所有测试
go test ./... -v

# 运行特定包的测试
go test ./internal/database -v
go test ./internal/handlers -v

# 运行集成测试
go test -v -run TestIntegration

# 查看测试覆盖率
go test ./... -cover
```

详细测试文档请查看 [TESTING.md](TESTING.md)

## 🎯 扩展建议

- [x] 添加单元测试 ✅
- [ ] 集成真实数据库（PostgreSQL/MySQL）
- [ ] 添加JWT认证
- [ ] 实现分页功能
- [ ] 添加日志系统
- [ ] 实现配置管理
- [ ] 添加Docker支持
- [ ] 实现API限流

## 📄 许可证

MIT License

## 🤝 贡献

欢迎提交Issue和Pull Request！

---

**Happy Coding! 🎉**

