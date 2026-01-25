package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"golang-web-demo/internal/database"
	"golang-web-demo/internal/models"

	"github.com/gin-gonic/gin"
)

// UserHandler 用户处理器
// 展示了Go语言的组合模式
// 将数据库依赖注入到处理器中
type UserHandler struct {
	db *database.InMemoryDB
}

// NewUserHandler 创建用户处理器
// 展示了依赖注入模式
func NewUserHandler(db *database.InMemoryDB) *UserHandler {
	return &UserHandler{db: db}
}

// CreateUser 创建用户
// 展示了：
// 1. HTTP请求处理
// 2. 参数绑定和验证
// 3. JSON响应
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest

	// 绑定JSON请求体到结构体
	// Gin会自动进行验证（基于binding标签）
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "请求参数无效",
			"details": err.Error(),
		})
		return
	}

	// 调用数据库层创建用户
	user, err := h.db.CreateUser(req.Name, req.Email, req.Age)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusCreated, gin.H{
		"message": "用户创建成功",
		"data":    user,
	})
}

// GetUser 获取单个用户
// 展示了：
// 1. URL参数获取
// 2. 类型转换（string -> int）
// 3. 错误处理
func (h *UserHandler) GetUser(c *gin.Context) {
	// 从URL路径获取参数
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的用户ID",
		})
		return
	}

	user, err := h.db.GetUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

// GetAllUsers 获取所有用户
// 展示了：
// 1. 列表响应
// 2. 空数组的处理
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users := h.db.GetAllUsers()

	// 如果用户列表为空，返回空数组而不是nil
	if users == nil {
		users = []*models.User{}
	}
	// resp := models.UserResponse{Count: len(users), Data: users}
	// c.JSON(http.StatusOK, resp)
	c.JSON(http.StatusOK, gin.H{
		"count": len(users),
		"data":  users,
	})
}

// UpdateUser 更新用户
// 展示了：
// 1. 部分更新模式
// 2. 指针的使用（区分零值和未提供值）
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的用户ID",
		})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "请求参数无效",
			"details": err.Error(),
		})
		return
	}

	// 将请求中的值转换为指针
	// 如果字段为零值，则传递nil（表示不更新该字段）
	var namePtr, emailPtr *string
	var agePtr *int

	if req.Name != "" {
		namePtr = &req.Name
	}
	if req.Email != "" {
		emailPtr = &req.Email
	}
	if req.Age > 0 {
		agePtr = &req.Age
	}

	user, err := h.db.UpdateUser(id, namePtr, emailPtr, agePtr)
	if err != nil {
		statusCode := http.StatusNotFound
		// 检查是否是邮箱冲突错误（通过错误消息判断）
		errMsg := err.Error()
		if strings.Contains(errMsg, "已被其他用户使用") {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, gin.H{
			"error": errMsg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "用户更新成功",
		"data":    user,
	})
}

// DeleteUser 删除用户
// 展示了：
// 1. DELETE操作的处理
// 2. 204 No Content响应
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的用户ID",
		})
		return
	}

	if err := h.db.DeleteUser(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 删除成功，返回204 No Content
	c.Status(http.StatusNoContent)
}
