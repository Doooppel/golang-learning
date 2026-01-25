package models

import (
	"time"
)

// User 用户模型
// 展示了Go语言的结构体定义
// 结构体是Go语言中组织数据的主要方式
type User struct {
	ID        int       `json:"id"`         // 用户ID
	Name      string    `json:"name"`       // 用户名
	Email     string    `json:"email"`      // 邮箱
	Age       int       `json:"age"`        // 年龄
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}

type UserResponse struct {
	Count int     `json:"count"`
	Data  []*User `json:"data"`
}

// CreateUserRequest 创建用户请求
// 展示了请求DTO（Data Transfer Object）的使用
type CreateUserRequest struct {
	Name  string `json:"name" binding:"required,min=2,max=50"` // 必填，2-50字符
	Email string `json:"email" binding:"required,email"`       // 必填，必须是邮箱格式
	Age   int    `json:"age" binding:"required,min=1,max=150"` // 必填，1-150之间
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Name  string `json:"name" binding:"omitempty,min=2,max=50"` // 可选，如果提供则2-50字符
	Email string `json:"email" binding:"omitempty,email"`       // 可选，如果提供则必须是邮箱格式
	Age   int    `json:"age" binding:"omitempty,min=1,max=150"` // 可选，如果提供则1-150之间
}
