package database

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"golang-web-demo/internal/models"
)

// InMemoryDB 内存数据库
// 展示了Go语言的高级特性：
// 1. 结构体封装
// 2. sync.RWMutex 读写锁（并发安全）
// 3. map数据结构的使用
// 4. 文件持久化
type InMemoryDB struct {
	users    map[int]*models.User // 使用map存储用户数据，key是用户ID
	mu       sync.RWMutex         // 读写锁，保证并发安全
	nextID   int                  // 下一个可用的ID
	filePath string               // 数据文件路径
}

// NewInMemoryDB 创建新的内存数据库实例
// 展示了Go语言的构造函数模式
// 会自动从文件加载数据（如果文件存在）
func NewInMemoryDB() *InMemoryDB {
	// 设置数据文件路径（在项目根目录下的data文件夹）
	filePath := "data/users.json"

	db := &InMemoryDB{
		users:    make(map[int]*models.User),
		nextID:   1,
		filePath: filePath,
	}

	// 尝试从文件加载数据
	if err := db.loadFromFile(); err != nil {
		// 如果文件不存在或加载失败，使用空数据库（这是正常的首次启动情况）
		fmt.Printf("⚠️  无法加载数据文件（首次启动或文件不存在）: %v\n", err)
	}

	return db
}

// CreateUser 创建用户
// 展示了：
// 1. 指针接收者方法
// 2. 错误处理
// 3. 时间处理
func (db *InMemoryDB) CreateUser(name, email string, age int) (*models.User, error) {
	db.mu.Lock()         // 写锁
	defer db.mu.Unlock() // defer确保函数返回时释放锁

	// 检查邮箱是否已存在
	for _, user := range db.users {
		if user.Email == email {
			return nil, fmt.Errorf("邮箱 %s 已存在", email)
		}
	}

	// 创建新用户
	now := time.Now()
	user := &models.User{
		ID:        db.nextID,
		Name:      name,
		Email:     email,
		Age:       age,
		CreatedAt: now,
		UpdatedAt: now,
	}

	db.users[user.ID] = user
	db.nextID++

	// 自动保存到文件
	if err := db.saveToFile(); err != nil {
		// 记录错误但不影响创建操作
		fmt.Printf("⚠️  保存数据到文件失败: %v\n", err)
	}

	return user, nil
}

// GetUser 根据ID获取用户
// 展示了：
// 1. 读锁的使用（允许多个读操作并发）
// 2. map的查找操作
func (db *InMemoryDB) GetUser(id int) (*models.User, error) {
	db.mu.RLock()         // 读锁
	defer db.mu.RUnlock() // defer释放读锁

	user, exists := db.users[id]
	if !exists {
		return nil, fmt.Errorf("用户 ID %d 不存在", id)
	}

	return user, nil
}

// GetAllUsers 获取所有用户
// 展示了：
// 1. slice的创建和追加
// 2. range遍历map
func (db *InMemoryDB) GetAllUsers() []*models.User {
	db.mu.RLock()
	defer db.mu.RUnlock()

	users := make([]*models.User, 0, len(db.users))
	for _, user := range db.users {
		users = append(users, user)
	}

	return users
}

// UpdateUser 更新用户
// 展示了：
// 1. 条件更新（只更新提供的字段）
// 2. 错误处理模式
func (db *InMemoryDB) UpdateUser(id int, name, email *string, age *int) (*models.User, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	user, exists := db.users[id]
	if !exists {
		return nil, fmt.Errorf("用户 ID %d 不存在", id)
	}

	// 如果提供了新值，则更新
	// 展示了指针的使用和nil检查
	if name != nil {
		user.Name = *name
	}
	if email != nil {
		// 检查邮箱是否被其他用户使用
		for _, u := range db.users {
			if u.ID != id && u.Email == *email {
				return nil, fmt.Errorf("邮箱 %s 已被其他用户使用", *email)
			}
		}
		user.Email = *email
	}
	if age != nil {
		user.Age = *age
	}

	user.UpdatedAt = time.Now()

	// 自动保存到文件
	if err := db.saveToFile(); err != nil {
		fmt.Printf("⚠️  保存数据到文件失败: %v\n", err)
	}

	return user, nil
}

// DeleteUser 删除用户
// 展示了：
// 1. map的删除操作
// 2. Go语言的delete内置函数
func (db *InMemoryDB) DeleteUser(id int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, exists := db.users[id]; !exists {
		return fmt.Errorf("用户 ID %d 不存在", id)
	}

	delete(db.users, id)

	// 自动保存到文件
	if err := db.saveToFile(); err != nil {
		fmt.Printf("⚠️  保存数据到文件失败: %v\n", err)
	}

	return nil
}

// GetUserByEmail 根据邮箱获取用户（辅助方法）
// 展示了：
// 1. 函数式编程风格
// 2. 线性搜索
func (db *InMemoryDB) GetUserByEmail(email string) (*models.User, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	for _, user := range db.users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, fmt.Errorf("邮箱 %s 不存在", email)
}

// saveToFile 保存数据到JSON文件
// 展示了：
// 1. JSON序列化
// 2. 文件操作
// 3. 目录创建
func (db *InMemoryDB) saveToFile() error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	// 创建数据目录（如果不存在）
	dir := filepath.Dir(db.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 将用户map转换为slice以便JSON序列化
	usersList := make([]*models.User, 0, len(db.users))
	for _, user := range db.users {
		usersList = append(usersList, user)
	}

	// 准备要保存的数据结构
	data := struct {
		Users  []*models.User `json:"users"`
		NextID int            `json:"next_id"`
	}{
		Users:  usersList,
		NextID: db.nextID,
	}

	// 序列化为JSON（格式化输出，便于阅读）
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("JSON序列化失败: %w", err)
	}

	// 写入文件（使用原子写入：先写临时文件，再重命名）
	tempFile := db.filePath + ".tmp"
	if err := os.WriteFile(tempFile, jsonData, 0644); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	// 原子性替换原文件
	if err := os.Rename(tempFile, db.filePath); err != nil {
		return fmt.Errorf("重命名文件失败: %w", err)
	}

	return nil
}

// loadFromFile 从JSON文件加载数据
// 展示了：
// 1. JSON反序列化
// 2. 文件读取
// 3. 错误处理
func (db *InMemoryDB) loadFromFile() error {
	// 检查文件是否存在
	if _, err := os.Stat(db.filePath); os.IsNotExist(err) {
		return fmt.Errorf("数据文件不存在")
	}

	// 读取文件内容
	jsonData, err := os.ReadFile(db.filePath)
	if err != nil {
		return fmt.Errorf("读取文件失败: %w", err)
	}

	// 解析JSON数据
	var data struct {
		Users  []*models.User `json:"users"`
		NextID int            `json:"next_id"`
	}

	if err := json.Unmarshal(jsonData, &data); err != nil {
		return fmt.Errorf("JSON解析失败: %w", err)
	}

	// 重建用户map
	db.mu.Lock()
	defer db.mu.Unlock()

	db.users = make(map[int]*models.User)
	for _, user := range data.Users {
		db.users[user.ID] = user
		// 更新nextID为最大ID+1
		if user.ID >= db.nextID {
			db.nextID = user.ID + 1
		}
	}

	// 如果文件中指定了nextID，使用它（更准确）
	if data.NextID > 0 {
		db.nextID = data.NextID
	}

	fmt.Printf("✅ 成功加载 %d 个用户数据\n", len(db.users))
	return nil
}
