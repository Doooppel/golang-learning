package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"golang-web-demo/internal/database"
	"golang-web-demo/internal/handlers"
	"golang-web-demo/internal/router"

	"github.com/gin-gonic/gin"
)

// TestIntegrationCRUD 集成测试：完整的CRUD流程
func TestIntegrationCRUD(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// 清理测试数据
	cleanupTestData()

	// 设置数据库和路由
	db := database.NewInMemoryDB()
	userHandler := handlers.NewUserHandler(db)
	r := router.SetupRouter(userHandler)

	// 1. 创建用户
	createReq := map[string]interface{}{
		"name":  "集成测试用户",
		"email": "integration@example.com",
		"age":   25,
	}
	jsonData, _ := json.Marshal(createReq)

	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("创建用户失败，状态码: %d, 响应: %s", w.Code, w.Body.String())
	}

	var createResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResp)
	userData := createResp["data"].(map[string]interface{})
	userID := int(userData["id"].(float64))

	// 2. 获取所有用户
	req2, _ := http.NewRequest("GET", "/api/v1/users", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Fatalf("获取所有用户失败，状态码: %d", w2.Code)
	}

	var listResp map[string]interface{}
	json.Unmarshal(w2.Body.Bytes(), &listResp)
	if listResp["count"].(float64) != 1 {
		t.Errorf("期望1个用户，实际是 %v", listResp["count"])
	}

	// 3. 获取单个用户
	req3, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/users/%d", userID), nil)
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)

	if w3.Code != http.StatusOK {
		t.Fatalf("获取用户失败，状态码: %d", w3.Code)
	}

	// 4. 更新用户
	updateReq := map[string]interface{}{
		"name": "集成测试用户更新",
		"age":  30,
	}
	updateData, _ := json.Marshal(updateReq)

	req4, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/users/%d", userID), bytes.NewBuffer(updateData))
	req4.Header.Set("Content-Type", "application/json")
	w4 := httptest.NewRecorder()
	r.ServeHTTP(w4, req4)

	if w4.Code != http.StatusOK {
		t.Fatalf("更新用户失败，状态码: %d", w4.Code)
	}

	// 5. 删除用户
	req5, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/users/%d", userID), nil)
	w5 := httptest.NewRecorder()
	r.ServeHTTP(w5, req5)

	if w5.Code != http.StatusNoContent {
		t.Fatalf("删除用户失败，状态码: %d", w5.Code)
	}

	// 6. 验证用户已删除
	req6, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/users/%d", userID), nil)
	w6 := httptest.NewRecorder()
	r.ServeHTTP(w6, req6)

	if w6.Code != http.StatusNotFound {
		t.Errorf("用户应该已被删除，状态码: %d", w6.Code)
	}
}

// TestIntegrationPersistence 集成测试：数据持久化
func TestIntegrationPersistence(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cleanupTestData()

	// 第一轮：创建数据
	db1 := database.NewInMemoryDB()
	userHandler1 := handlers.NewUserHandler(db1)
	r1 := router.SetupRouter(userHandler1)

	createReq := map[string]interface{}{
		"name":  "持久化测试",
		"email": "persist@example.com",
		"age":   25,
	}
	jsonData, _ := json.Marshal(createReq)

	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r1.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("创建用户失败，状态码: %d", w.Code)
	}

	// 第二轮：从新数据库实例加载数据
	db2 := database.NewInMemoryDB()
	userHandler2 := handlers.NewUserHandler(db2)
	r2 := router.SetupRouter(userHandler2)

	req2, _ := http.NewRequest("GET", "/api/v1/users", nil)
	w2 := httptest.NewRecorder()
	r2.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Fatalf("获取用户失败，状态码: %d", w2.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w2.Body.Bytes(), &resp)
	if resp["count"].(float64) != 1 {
		t.Errorf("应该加载1个用户，实际是 %v", resp["count"])
	}
}

// TestIntegrationHealthCheck 集成测试：健康检查
func TestIntegrationHealthCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	db := database.NewInMemoryDB()
	userHandler := handlers.NewUserHandler(db)
	r := router.SetupRouter(userHandler)

	req, _ := http.NewRequest("GET", "/api/v1/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("健康检查失败，状态码: %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["status"] != "ok" {
		t.Errorf("期望status为'ok'，实际是 '%v'", resp["status"])
	}
}

// cleanupTestData 清理测试数据
func cleanupTestData() {
	testFile := "data/users.json"
	os.Remove(testFile)
	os.Remove(testFile + ".tmp")
}

