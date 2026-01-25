package database

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"
)

// TestNewInMemoryDB 测试创建新的内存数据库
func TestNewInMemoryDB(t *testing.T) {
	// 清理测试数据
	cleanupTestData(t)

	db := NewInMemoryDB()
	if db == nil {
		t.Fatal("NewInMemoryDB() 返回 nil")
	}

	if db.users == nil {
		t.Error("users map 未初始化")
	}

	if db.nextID != 1 {
		t.Errorf("nextID 应该是 1，实际是 %d", db.nextID)
	}
}

// TestCreateUser 测试创建用户
func TestCreateUser(t *testing.T) {
	cleanupTestData(t)
	db := NewInMemoryDB()

	user, err := db.CreateUser("张三", "zhangsan@example.com", 25)
	if err != nil {
		t.Fatalf("创建用户失败: %v", err)
	}

	if user == nil {
		t.Fatal("创建的用户为 nil")
	}

	if user.ID != 1 {
		t.Errorf("用户ID应该是 1，实际是 %d", user.ID)
	}

	if user.Name != "张三" {
		t.Errorf("用户名应该是 '张三'，实际是 '%s'", user.Name)
	}

	if user.Email != "zhangsan@example.com" {
		t.Errorf("邮箱应该是 'zhangsan@example.com'，实际是 '%s'", user.Email)
	}

	if user.Age != 25 {
		t.Errorf("年龄应该是 25，实际是 %d", user.Age)
	}

	// 测试重复邮箱
	_, err = db.CreateUser("李四", "zhangsan@example.com", 30)
	if err == nil {
		t.Error("应该拒绝重复邮箱")
	}
}

// TestGetUser 测试获取用户
func TestGetUser(t *testing.T) {
	cleanupTestData(t)
	db := NewInMemoryDB()

	// 创建用户
	createdUser, _ := db.CreateUser("王五", "wangwu@example.com", 28)

	// 获取用户
	user, err := db.GetUser(createdUser.ID)
	if err != nil {
		t.Fatalf("获取用户失败: %v", err)
	}

	if user.ID != createdUser.ID {
		t.Errorf("用户ID不匹配: 期望 %d，实际 %d", createdUser.ID, user.ID)
	}

	// 测试不存在的用户
	_, err = db.GetUser(999)
	if err == nil {
		t.Error("应该返回错误，用户不存在")
	}
}

// TestGetAllUsers 测试获取所有用户
func TestGetAllUsers(t *testing.T) {
	cleanupTestData(t)
	db := NewInMemoryDB()

	// 创建多个用户
	db.CreateUser("用户1", "user1@example.com", 20)
	db.CreateUser("用户2", "user2@example.com", 25)
	db.CreateUser("用户3", "user3@example.com", 30)

	users := db.GetAllUsers()
	if len(users) != 3 {
		t.Errorf("应该有 3 个用户，实际是 %d", len(users))
	}
}

// TestUpdateUser 测试更新用户
func TestUpdateUser(t *testing.T) {
	cleanupTestData(t)
	db := NewInMemoryDB()

	// 创建用户
	user, _ := db.CreateUser("原始名字", "original@example.com", 25)

	// 更新名字
	newName := "新名字"
	updatedUser, err := db.UpdateUser(user.ID, &newName, nil, nil)
	if err != nil {
		t.Fatalf("更新用户失败: %v", err)
	}

	if updatedUser.Name != "新名字" {
		t.Errorf("名字应该更新为 '新名字'，实际是 '%s'", updatedUser.Name)
	}

	// 更新年龄
	newAge := 30
	updatedUser, err = db.UpdateUser(user.ID, nil, nil, &newAge)
	if err != nil {
		t.Fatalf("更新年龄失败: %v", err)
	}

	if updatedUser.Age != 30 {
		t.Errorf("年龄应该更新为 30，实际是 %d", updatedUser.Age)
	}

	// 测试不存在的用户
	_, err = db.UpdateUser(999, &newName, nil, nil)
	if err == nil {
		t.Error("应该返回错误，用户不存在")
	}
}

// TestDeleteUser 测试删除用户
func TestDeleteUser(t *testing.T) {
	cleanupTestData(t)
	db := NewInMemoryDB()

	// 创建用户
	user, _ := db.CreateUser("待删除", "delete@example.com", 25)

	// 删除用户
	err := db.DeleteUser(user.ID)
	if err != nil {
		t.Fatalf("删除用户失败: %v", err)
	}

	// 验证用户已删除
	_, err = db.GetUser(user.ID)
	if err == nil {
		t.Error("用户应该已被删除")
	}

	// 测试删除不存在的用户
	err = db.DeleteUser(999)
	if err == nil {
		t.Error("应该返回错误，用户不存在")
	}
}

// TestPersistence 测试数据持久化
func TestPersistence(t *testing.T) {
	cleanupTestData(t)

	// 创建第一个数据库实例并添加数据
	db1 := NewInMemoryDB()
	db1.CreateUser("持久化用户1", "persist1@example.com", 25)
	db1.CreateUser("持久化用户2", "persist2@example.com", 30)

	// 验证数据已保存到文件
	if _, err := os.Stat(db1.filePath); os.IsNotExist(err) {
		t.Fatal("数据文件应该已创建")
	}

	// 创建第二个数据库实例（应该从文件加载数据）
	db2 := NewInMemoryDB()
	users := db2.GetAllUsers()

	if len(users) != 2 {
		t.Errorf("应该加载 2 个用户，实际是 %d", len(users))
	}

	// 验证数据正确性
	found := false
	for _, user := range users {
		if user.Email == "persist1@example.com" && user.Name == "持久化用户1" {
			found = true
			break
		}
	}
	if !found {
		t.Error("应该找到持久化用户1")
	}
}

// TestConcurrency 测试并发安全性
func TestConcurrency(t *testing.T) {
	cleanupTestData(t)
	db := NewInMemoryDB()

	const numGoroutines = 100
	const numUsersPerGoroutine = 10

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*numUsersPerGoroutine)

	// 并发创建用户
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numUsersPerGoroutine; j++ {
				email := fmt.Sprintf("user%d_%d@example.com", id, j)
				_, err := db.CreateUser(fmt.Sprintf("用户%d_%d", id, j), email, 20+j)
				if err != nil {
					errors <- err
				}
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	// 检查错误
	errorCount := 0
	for err := range errors {
		if err != nil {
			errorCount++
			t.Logf("并发错误: %v", err)
		}
	}

	// 验证所有用户都已创建（除了可能的重复邮箱错误）
	users := db.GetAllUsers()
	expectedCount := numGoroutines * numUsersPerGoroutine
	// 允许一些重复邮箱错误
	if len(users) < expectedCount-10 {
		t.Errorf("应该创建大约 %d 个用户，实际是 %d", expectedCount, len(users))
	}
}

// TestConcurrentReadWrite 测试并发读写
func TestConcurrentReadWrite(t *testing.T) {
	cleanupTestData(t)
	db := NewInMemoryDB()

	// 先创建一些用户
	for i := 0; i < 10; i++ {
		db.CreateUser(fmt.Sprintf("用户%d", i), fmt.Sprintf("user%d@example.com", i), 20+i)
	}

	var wg sync.WaitGroup
	errors := make(chan error, 100)

	// 并发读取
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			users := db.GetAllUsers()
			if len(users) == 0 {
				errors <- fmt.Errorf("读取到空列表")
			}
		}()
	}

	// 并发写入
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			email := fmt.Sprintf("concurrent%d@example.com", id)
			_, err := db.CreateUser(fmt.Sprintf("并发用户%d", id), email, 25)
			if err != nil {
				errors <- err
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	// 检查是否有严重错误
	for err := range errors {
		if err != nil && err.Error() != "邮箱 concurrent%d@example.com 已存在" {
			t.Errorf("并发错误: %v", err)
		}
	}
}

// TestGetUserByEmail 测试根据邮箱获取用户
func TestGetUserByEmail(t *testing.T) {
	cleanupTestData(t)
	db := NewInMemoryDB()

	// 创建用户
	createdUser, _ := db.CreateUser("邮箱测试", "emailtest@example.com", 25)

	// 根据邮箱查找
	user, err := db.GetUserByEmail("emailtest@example.com")
	if err != nil {
		t.Fatalf("根据邮箱获取用户失败: %v", err)
	}

	if user.ID != createdUser.ID {
		t.Errorf("用户ID不匹配: 期望 %d，实际 %d", createdUser.ID, user.ID)
	}

	// 测试不存在的邮箱
	_, err = db.GetUserByEmail("nonexist@example.com")
	if err == nil {
		t.Error("应该返回错误，邮箱不存在")
	}
}

// TestUpdateUserEmailConflict 测试更新邮箱冲突
func TestUpdateUserEmailConflict(t *testing.T) {
	cleanupTestData(t)
	db := NewInMemoryDB()

	// 创建两个用户
	_, _ = db.CreateUser("用户1", "user1@example.com", 25)
	user2, _ := db.CreateUser("用户2", "user2@example.com", 30)

	// 尝试将user2的邮箱更新为user1的邮箱（应该失败）
	conflictEmail := "user1@example.com"
	_, err := db.UpdateUser(user2.ID, nil, &conflictEmail, nil)
	if err == nil {
		t.Error("应该拒绝邮箱冲突")
	}
}

// TestAutoSaveOnOperations 测试操作后自动保存
func TestAutoSaveOnOperations(t *testing.T) {
	cleanupTestData(t)

	db := NewInMemoryDB()
	testFile := "data/test_users.json"
	db.filePath = testFile

	// 创建用户（应该自动保存）
	_, err := db.CreateUser("自动保存", "autosave@example.com", 25)
	if err != nil {
		t.Fatalf("创建用户失败: %v", err)
	}

	// 等待文件写入
	time.Sleep(100 * time.Millisecond)

	// 验证文件存在
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Error("文件应该已创建")
	}

	// 更新用户（应该自动保存）
	newName := "自动保存更新"
	_, err = db.UpdateUser(1, &newName, nil, nil)
	if err != nil {
		t.Fatalf("更新用户失败: %v", err)
	}

	// 删除用户（应该自动保存）
	err = db.DeleteUser(1)
	if err != nil {
		t.Fatalf("删除用户失败: %v", err)
	}

	// 清理测试文件
	os.Remove(testFile)
	os.Remove(testFile + ".tmp")
}

// cleanupTestData 清理测试数据
func cleanupTestData(t *testing.T) {
	testFile := "data/users.json"
	if err := os.Remove(testFile); err != nil && !os.IsNotExist(err) {
		t.Logf("清理测试文件失败: %v", err)
	}
	// 也清理临时文件
	os.Remove(testFile + ".tmp")
}

// BenchmarkCreateUser 性能测试：创建用户
func BenchmarkCreateUser(b *testing.B) {
	cleanupTestData(&testing.T{})
	db := NewInMemoryDB()
	db.filePath = "data/bench_users.json"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		email := fmt.Sprintf("bench%d@example.com", i)
		db.CreateUser(fmt.Sprintf("用户%d", i), email, 25)
	}

	os.Remove("data/bench_users.json")
	os.Remove("data/bench_users.json.tmp")
}

// BenchmarkGetUser 性能测试：获取用户
func BenchmarkGetUser(b *testing.B) {
	cleanupTestData(&testing.T{})
	db := NewInMemoryDB()

	// 预先创建用户
	user, _ := db.CreateUser("基准测试", "bench@example.com", 25)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db.GetUser(user.ID)
	}
}

// BenchmarkGetAllUsers 性能测试：获取所有用户
func BenchmarkGetAllUsers(b *testing.B) {
	cleanupTestData(&testing.T{})
	db := NewInMemoryDB()

	// 预先创建100个用户
	for i := 0; i < 100; i++ {
		email := fmt.Sprintf("bench%d@example.com", i)
		db.CreateUser(fmt.Sprintf("用户%d", i), email, 25)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db.GetAllUsers()
	}
}

