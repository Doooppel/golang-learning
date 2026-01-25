package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"golang-web-demo/internal/database"
	"golang-web-demo/internal/models"

	"github.com/gin-gonic/gin"
)

// setupTestRouter 设置测试路由
func setupTestRouter() (*gin.Engine, *UserHandler) {
	gin.SetMode(gin.TestMode)
	
	// 使用独立的测试数据文件，并确保清理
	testFile := "data/test_users.json"
	os.Remove(testFile)
	os.Remove(testFile + ".tmp")
	
	// 创建用于测试的数据库实例（不加载文件）
	db := database.NewInMemoryDBForTest(testFile)
	
	handler := NewUserHandler(db)
	router := gin.New()
	
	api := router.Group("/api/v1")
	{
		users := api.Group("/users")
		{
			users.POST("", handler.CreateUser)
			users.GET("", handler.GetAllUsers)
			users.GET("/:id", handler.GetUser)
			users.PUT("/:id", handler.UpdateUser)
			users.DELETE("/:id", handler.DeleteUser)
		}
	}
	
	return router, handler
}

// TestCreateUser 测试创建用户API
func TestCreateUser(t *testing.T) {
	router, _ := setupTestRouter()

	// 测试成功创建
	reqBody := models.CreateUserRequest{
		Name:  "测试用户",
		Email: "test@example.com",
		Age:   25,
	}
	jsonData, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("期望状态码 %d，实际是 %d", http.StatusCreated, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["message"] != "用户创建成功" {
		t.Errorf("期望消息 '用户创建成功'，实际是 '%v'", response["message"])
	}

	// 测试无效请求
	invalidReq, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer([]byte("invalid json")))
	invalidReq.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, invalidReq)

	if w2.Code != http.StatusBadRequest {
		t.Errorf("期望状态码 %d，实际是 %d", http.StatusBadRequest, w2.Code)
	}
}

// TestGetUser 测试获取用户API
func TestGetUser(t *testing.T) {
	router, handler := setupTestRouter()

	// 先创建一个用户
	_, _ = handler.db.CreateUser("获取测试", "gettest@example.com", 25)

	// 测试成功获取
	req, _ := http.NewRequest("GET", "/api/v1/users/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("期望状态码 %d，实际是 %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["data"] == nil {
		t.Error("响应应该包含data字段")
	}

	// 测试不存在的用户
	req2, _ := http.NewRequest("GET", "/api/v1/users/999", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusNotFound {
		t.Errorf("期望状态码 %d，实际是 %d", http.StatusNotFound, w2.Code)
	}

	// 测试无效ID
	req3, _ := http.NewRequest("GET", "/api/v1/users/invalid", nil)
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)

	if w3.Code != http.StatusBadRequest {
		t.Errorf("期望状态码 %d，实际是 %d", http.StatusBadRequest, w3.Code)
	}
}

// TestGetAllUsers 测试获取所有用户API
func TestGetAllUsers(t *testing.T) {
	router, handler := setupTestRouter()

	// 创建几个用户
	handler.db.CreateUser("用户1", "user1@example.com", 25)
	handler.db.CreateUser("用户2", "user2@example.com", 30)

	req, _ := http.NewRequest("GET", "/api/v1/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("期望状态码 %d，实际是 %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if count, ok := response["count"].(float64); !ok || count != 2 {
		t.Errorf("期望count为2，实际是 %v", response["count"])
	}
}

// TestUpdateUser 测试更新用户API
func TestUpdateUser(t *testing.T) {
	router, handler := setupTestRouter()

	// 先创建一个用户
	handler.db.CreateUser("原始名字", "original@example.com", 25)

	// 测试成功更新
	updateReq := models.UpdateUserRequest{
		Name: "新名字",
		Age:  30,
	}
	jsonData, _ := json.Marshal(updateReq)

	req, _ := http.NewRequest("PUT", "/api/v1/users/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("期望状态码 %d，实际是 %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["message"] != "用户更新成功" {
		t.Errorf("期望消息 '用户更新成功'，实际是 '%v'", response["message"])
	}

	// 验证更新后的数据
	user, _ := handler.db.GetUser(1)
	if user.Name != "新名字" || user.Age != 30 {
		t.Error("用户数据未正确更新")
	}
}

// TestDeleteUser 测试删除用户API
func TestDeleteUser(t *testing.T) {
	router, handler := setupTestRouter()

	// 先创建一个用户
	handler.db.CreateUser("待删除", "delete@example.com", 25)

	// 测试成功删除
	req, _ := http.NewRequest("DELETE", "/api/v1/users/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("期望状态码 %d，实际是 %d", http.StatusNoContent, w.Code)
	}

	// 验证用户已删除
	_, err := handler.db.GetUser(1)
	if err == nil {
		t.Error("用户应该已被删除")
	}
}

// TestCreateUserValidation 测试创建用户的验证
func TestCreateUserValidation(t *testing.T) {
	router, _ := setupTestRouter()

	testCases := []struct {
		name        string
		requestBody models.CreateUserRequest
		expectCode  int
	}{
		{
			name: "缺少必填字段",
			requestBody: models.CreateUserRequest{
				Name: "",
				Email: "test@example.com",
				Age: 25,
			},
			expectCode: http.StatusBadRequest,
		},
		{
			name: "无效邮箱",
			requestBody: models.CreateUserRequest{
				Name: "测试",
				Email: "invalid-email",
				Age: 25,
			},
			expectCode: http.StatusBadRequest,
		},
		{
			name: "年龄超出范围",
			requestBody: models.CreateUserRequest{
				Name: "测试",
				Email: "test@example.com",
				Age: 200,
			},
			expectCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.expectCode {
				t.Errorf("期望状态码 %d，实际是 %d", tc.expectCode, w.Code)
			}
		})
	}
}

