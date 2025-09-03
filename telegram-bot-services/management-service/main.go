/*
 * 文件功能描述：管理服务主程序，处理搜索API请求
 * 主要类/接口说明：主函数和HTTP处理函数
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
    
    "management-service/service"
)

// 全局服务实例
var searchService service.SearchService

// 初始化服务
// @author fcj
// @date 2023-11-15
// @version 1.0.0
func initServices() {
    // 初始化搜索服务
    searchConfig := service.SearchConfig{
        MeilisearchURL: "http://your-meilisearch-url",
        MeilisearchKey: "YOUR_MEILISEARCH_KEY",
    }
    searchService = service.NewSearchService(searchConfig)
}

// searchHandler 处理搜索请求
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param w 响应写入器
// @param r 请求
func searchHandler(w http.ResponseWriter, r *http.Request) {
    // 获取查询参数
    query := r.URL.Query().Get("q")
    page := r.URL.Query().Get("page")
    limit := r.URL.Query().Get("limit")
    filter := r.URL.Query().Get("filter")

    // 执行搜索
    result, err := searchService.Search(query, page, limit, filter)
    if err != nil {
        http.Error(w, "Search failed: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // 返回搜索结果
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}

func main() {
    // 初始化服务
    initServices()

    // 设置路由
    r := mux.NewRouter()
    r.HandleFunc("/api/search", searchHandler).Methods("GET")

    // 启动服务器
    log.Println("Management Service running on :8080")
    err := http.ListenAndServe(":8080", r)
    if err != nil {
        log.Fatal("Failed to start server: ", err)
    }
}

/*
 * 关键算法说明：
 * 1. 服务初始化：使用依赖注入模式初始化搜索服务
 * 2. API处理：接收HTTP请求并调用服务层处理
 * 
 * 待优化事项：
 * 1. 配置管理：从环境变量或配置文件加载配置
 * 2. 错误处理：改进错误处理和恢复机制
 * 3. 日志记录：添加结构化日志记录
 * 
 * 兼容性说明：
 * 1. 依赖service包中的服务接口和实现
 */