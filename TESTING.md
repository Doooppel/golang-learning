# 🧪 测试文档

本项目包含完整的测试套件，用于验证所有功能是否正常工作。

## 📋 测试文件

### 1. 数据库层测试
- **文件**: `internal/database/memory_db_test.go`
- **覆盖**: 所有数据库操作、持久化、并发安全

### 2. 处理器层测试
- **文件**: `internal/handlers/user_handler_test.go`
- **覆盖**: 所有HTTP API端点、参数验证

### 3. 集成测试
- **文件**: `test_integration_test.go`
- **覆盖**: 完整的CRUD流程、数据持久化、健康检查

## 🚀 运行测试

### 运行所有测试
```bash
go test ./... -v
```

### 运行特定包的测试
```bash
# 只测试数据库层
go test ./internal/database -v

# 只测试处理器层
go test ./internal/handlers -v

# 只运行集成测试
go test -v -run TestIntegration
```

### 运行性能测试（Benchmark）
```bash
go test ./internal/database -bench=.
```

### 查看测试覆盖率
```bash
go test ./... -cover
```

## 📊 测试覆盖范围

### 数据库层 (`internal/database`)
- ✅ 创建数据库实例
- ✅ 创建用户（包括重复邮箱检查）
- ✅ 获取用户（单个和全部）
- ✅ 更新用户（部分更新）
- ✅ 删除用户
- ✅ 根据邮箱查找用户
- ✅ 数据持久化（保存和加载）
- ✅ 并发安全性（100个goroutine并发操作）
- ✅ 并发读写测试

### 处理器层 (`internal/handlers`)
- ✅ 创建用户API
- ✅ 获取用户API（单个和全部）
- ✅ 更新用户API
- ✅ 删除用户API
- ✅ 参数验证（邮箱格式、年龄范围等）
- ✅ 错误处理

### 集成测试
- ✅ 完整的CRUD流程
- ✅ 数据持久化验证
- ✅ 健康检查端点

## 🔍 测试说明

### 测试隔离
每个测试都使用独立的测试数据文件（`data/test_users.json`），确保测试之间不会相互影响。

### 并发测试
包含专门的并发测试，验证在高并发场景下的数据安全性：
- `TestConcurrency`: 100个goroutine同时创建用户
- `TestConcurrentReadWrite`: 并发读写测试

### 性能测试
包含基准测试（Benchmark），用于评估性能：
- `BenchmarkCreateUser`: 创建用户性能
- `BenchmarkGetUser`: 获取用户性能
- `BenchmarkGetAllUsers`: 获取所有用户性能

## 📝 添加新测试

### 添加数据库测试
在 `internal/database/memory_db_test.go` 中添加新测试函数：
```go
func TestNewFeature(t *testing.T) {
    cleanupTestData(t)
    db := NewInMemoryDB()
    // 测试代码...
}
```

### 添加处理器测试
在 `internal/handlers/user_handler_test.go` 中添加新测试函数：
```go
func TestNewAPI(t *testing.T) {
    router, handler := setupTestRouter()
    // 测试代码...
}
```

### 添加集成测试
在 `test_integration_test.go` 中添加新测试函数：
```go
func TestIntegrationNewFeature(t *testing.T) {
    gin.SetMode(gin.TestMode)
    cleanupTestData()
    // 测试代码...
}
```

## ⚠️ 注意事项

1. **测试数据清理**: 测试会自动清理测试数据文件，但如果有测试失败，可能需要手动清理 `data/test_users.json`
2. **并发测试**: 并发测试可能需要几秒钟，这是正常的
3. **文件路径**: 测试使用独立的文件路径，不会影响生产数据

## 🎯 测试最佳实践

1. **每个测试应该独立**: 不依赖其他测试的执行顺序
2. **清理测试数据**: 使用 `cleanupTestData()` 确保测试环境干净
3. **使用表驱动测试**: 对于多个相似测试用例，使用表驱动测试
4. **测试错误情况**: 不仅要测试成功情况，也要测试错误情况
5. **并发测试**: 对于涉及共享状态的代码，必须进行并发测试

## 📈 持续集成

这些测试可以集成到CI/CD流程中：
```yaml
# 示例 GitHub Actions
- name: Run tests
  run: go test ./... -v -coverprofile=coverage.out
```

---

**Happy Testing! 🎉**

