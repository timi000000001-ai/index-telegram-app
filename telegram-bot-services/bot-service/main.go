/*
 * 文件功能描述：机器人服务主程序，处理Telegram机器人的初始化和Webhook处理
 * 主要类/接口说明：主函数和Webhook处理函数
 * 修改历史记录：
 * @author fcj
 * @date 2023-11-15
 * @version 1.0.0
 * © Telegram Bot Services Team
 */

package main

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/gorilla/mux"
    "gopkg.in/telebot.v3"
    
    "bot-service/service"
)

// 全局服务实例
var (
    botService     service.BotService
    messageService service.MessageService
    storageService service.StorageService
    searchService  service.SearchService
)

// Webhook 路由处理
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param w 响应写入器
// @param r 请求
func webhookHandler(w http.ResponseWriter, r *http.Request) {
    token := mux.Vars(r)["token"]
    
    update := &telebot.Update{}
    if err := json.NewDecoder(r.Body).Decode(update); err != nil {
        http.Error(w, "Invalid update", http.StatusBadRequest)
        return
    }

    if err := botService.ProcessUpdate(token, update); err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    
    w.WriteHeader(http.StatusOK)
}

// 初始化服务
// @author fcj
// @date 2023-11-15
// @version 1.0.0
func initServices() {
    // 初始化存储服务
    storageConfig := service.StorageConfig{
        PocketBaseURL:    "http://your-pocketbase-url",
        MeilisearchURL:   "http://your-meilisearch-url",
        MeilisearchToken: "YOUR_MEILISEARCH_TOKEN",
    }
    storageService = service.NewStorageService(storageConfig)
    
    // 初始化搜索服务
    searchConfig := service.SearchConfig{
        MeilisearchURL:      "http://your-meilisearch-url",
        MeilisearchKey:      "YOUR_MEILISEARCH_TOKEN",
        ManagementServiceURL: "http://your-management-service:8080",
    }
    searchService = service.NewSearchService(searchConfig)
    
    // 初始化消息服务
    messageService = service.NewMessageService(storageService, searchService)
    
    // 初始化机器人服务
    botService = service.NewBotService()
}

// 初始化机器人
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param configs 机器人配置列表
func initBots(configs []service.BotConfig) {
    for _, config := range configs {
        if err := botService.InitBot(config); err != nil {
            log.Printf("Failed to init bot %s: %v", config.Token, err)
            continue
        }
        
        // 获取机器人实例并注册处理函数
        bot, exists := botService.GetBot(config.Token)
        if exists {
            botService.RegisterHandlers(bot, messageService)
        }
    }
}

func main() {
    // 初始化服务
    initServices()
    
    // 机器人配置
    botConfigs := []service.BotConfig{
        {Token: "YOUR_BOT_TOKEN_1", WebhookURL: "https://your-bot-service.com/webhook/YOUR_BOT_TOKEN_1"},
        {Token: "YOUR_BOT_TOKEN_2", WebhookURL: "https://your-bot-service.com/webhook/YOUR_BOT_TOKEN_2"},
    }
    
    // 初始化机器人
    initBots(botConfigs)

    // 设置路由
    r := mux.NewRouter()
    r.HandleFunc("/webhook/{token}", webhookHandler).Methods("POST")

    // 启动服务器
    log.Println("Bot Service running on :8081")
    err := http.ListenAndServeTLS(":8081", "cert.pem", "key.pem", r)
    if err != nil {
        log.Fatal("Failed to start server: ", err)
    }
}

/*
 * 关键算法说明：
 * 1. 服务初始化：使用依赖注入模式初始化各个服务
 * 2. Webhook处理：接收和处理Telegram Webhook更新
 * 
 * 待优化事项：
 * 1. 配置管理：从环境变量或配置文件加载配置
 * 2. 错误处理：改进错误处理和恢复机制
 * 3. 日志记录：添加结构化日志记录
 * 
 * 兼容性说明：
 * 1. 依赖service包中的服务接口和实现
 * 2. 需要有效的TLS证书（cert.pem和key.pem）
 */