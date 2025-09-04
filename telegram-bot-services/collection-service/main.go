/*
 * 文件功能描述：采集服务主程序，负责从Telegram群组采集消息
 * 主要类/接口说明：主函数及HTTP路由处理
 * 修改历史记录：
 * @author fcj
 * @date 2023-11-15
 * @version 1.0.0
 * © Telegram Bot Services Team
 */

package main

import (
	"collection-service/internal/config"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"os"
	"time"

	"collection-service/service"
)

// ConfigureRequest 配置请求结构体
// @author fcj
// @date 2023-11-15
// @version 1.0.0
type ConfigureRequest struct {
	ChatIDs []int64 `json:"chat_ids"`
}

// LoginRequest 登录请求结构体
// @author fcj
// @date 2023-11-15
// @version 1.0.0
type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Code        string `json:"code,omitempty"`
}

// 全局服务实例
var (
	telegramService  service.TelegramService
	sessionService   service.SessionService
	storageService   service.StorageService
	collectionService service.CollectionService
)

// handleLogin 处理登录请求
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param w http.ResponseWriter
// @param r *http.Request
func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	if req.PhoneNumber == "" {
		http.Error(w, "Phone number is required", http.StatusBadRequest)
		return
	}

	// 如果没有提供验证码，则发送验证码
	if req.Code == "" {
		// 创建会话
		if err := sessionService.CreateSession(req.PhoneNumber); err != nil {
			log.Printf("Failed to create session: %v", err)
		}

		// 发送验证码
		if err := telegramService.SendCode(req.PhoneNumber); err != nil {
			http.Error(w, fmt.Sprintf("Failed to send verification code: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "code_sent"})
		return
	}

	// 如果提供了验证码，则进行登录
	if err := telegramService.Login(req.PhoneNumber, req.Code); err != nil {
		http.Error(w, fmt.Sprintf("Login failed: %v", err), http.StatusInternalServerError)
		return
	}

	// 检查是否已认证
	isAuthorized, err := telegramService.IsAuthorized(req.PhoneNumber)
	if err != nil || !isAuthorized {
		http.Error(w, fmt.Sprintf("Authorization check failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "logged_in"})
}

// handleConfigure 处理配置请求
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param w http.ResponseWriter
// @param r *http.Request
func handleConfigure(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ConfigureRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	// 获取请求中的电话号码
	phoneNumber := r.URL.Query().Get("phone_number")
	if phoneNumber == "" {
		http.Error(w, "Phone number is required", http.StatusBadRequest)
		return
	}

	// 更新会话配置
	if err := sessionService.UpdateSession(phoneNumber, req.ChatIDs); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update session: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "configured"})
}

// handleStartCollection 处理开始采集请求
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param w http.ResponseWriter
// @param r *http.Request
func handleStartCollection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 开始采集
	if err := collectionService.StartCollection(); err != nil {
		http.Error(w, fmt.Sprintf("Failed to start collection: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "collection_started"})
}

// handleHealth 处理健康检查请求
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param w http.ResponseWriter
// @param r *http.Request
func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

// initServices 初始化服务
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @return error 错误信息
func initServices(cfg *config.Config) error {
	// 初始化服务

	storageService = service.NewStorageService(service.StorageConfig{PocketBaseURL: cfg.Storage.PocketBaseURL})
	sessionService = service.NewSessionService()
	telegramService = service.NewTelegramService(service.TelegramConfig{
		AppID:   cfg.Telegram.AppID,
		AppHash: cfg.Telegram.AppHash,
	}, storageService)
	collectionService = service.NewCollectionService(service.SearchConfig{
		MeilisearchURL:  cfg.Search.MeilisearchURL,
		MeilisearchToken: cfg.Search.MeilisearchToken,
		MessageLimit:    cfg.Search.MessageLimit,
	}, telegramService, sessionService)

	return nil
}

// main 主函数
// @author fcj
// @date 2023-11-15
// @version 1.0.0
func main() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	cfg, err := config.Load(env)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}



	// 初始化服务
	if err := initServices(cfg); err != nil {
		log.Fatalf("Failed to initialize services: %v", err)
	}

	// 设置路由
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/configure", handleConfigure)
	http.HandleFunc("/collect", handleStartCollection)
	http.HandleFunc("/health", handleHealth)

	// 启动定时采集
	if err := collectionService.ScheduleCollection(5 * time.Minute); err != nil {
		log.Printf("Warning: Failed to schedule collection: %v", err)
	}

	// 启动HTTP服务器
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Starting collection service on port %s", cfg.Server.Port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

/*
 * 关键算法说明：
 * 1. 服务初始化：使用依赖注入模式初始化各个服务
 * 2. HTTP路由：处理登录、配置和采集请求
 * 3. 定时采集：启动定时任务定期采集消息
 * 
 * 待优化事项：
 * 1. 路由框架：使用更强大的路由框架如Gin或Echo
 * 2. 中间件：添加日志、认证和错误处理中间件
 * 3. 配置管理：使用结构化配置管理
 * 
 * 兼容性说明：
 * 1. 依赖环境变量配置
 * 2. 使用HTTP API与其他服务交互
 */