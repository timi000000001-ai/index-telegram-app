/*
 * 文件功能描述：管理服务主程序，处理搜索API请求并集成PocketBase
 * 主要类/接口说明：主函数和HTTP处理函数
 * 修改历史记录：
 * @author fcj
 * @date 2023-11-15
 * @version 1.2.0
 * © Telegram Bot Services Team
 */

package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"management-service/internal/config"
	_ "management-service/migrations"
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
	cfg := config.Get()
	searchConfig := service.SearchConfig{
		MeilisearchURL: cfg.MeilisearchURL,
		MeilisearchKey: cfg.MeilisearchKey,
	}
	searchService = service.NewSearchService(searchConfig)
}

// searchHandler 处理搜索请求 (适配 PocketBase v0.29.x)
// @author fcj
// @date 2024-07-26
// @version 1.1.0
// @param e PocketBase 请求事件
func searchHandler(e *core.RequestEvent) error {
	// 获取查询参数
	query := e.Request.URL.Query().Get("q")
	page := e.Request.URL.Query().Get("page")
	limit := e.Request.URL.Query().Get("limit")
	filter := e.Request.URL.Query().Get("filter")

	// 执行搜索
	result, err := searchService.Search(query, page, limit, filter)
	if err != nil {
		return err // PocketBase 会处理错误并返回 500
	}

	// 返回搜索结果
	return e.JSON(http.StatusOK, result)
}

func main() {
	// 加载配置
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}
	configPath := filepath.Join(wd, "configs")
	if _, err := config.Load(configPath); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 初始化服务
	initServices()

	// --- PocketBase ---
	// 获取可执行文件的路径
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}
	dataDir := filepath.Join(filepath.Dir(exePath), "pb_data")

	app := pocketbase.NewWithConfig(pocketbase.Config{
		DefaultDataDir: dataDir,
	})

	// 将现有的路由集成到 PocketBase 的服务器中
	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		// 注册 /api/search 路由
		e.Router.GET("/api/search", searchHandler)
		return e.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

/*
 * 关键算法说明：
 * 1. 服务初始化：使用依赖注入模式初始化搜索服务
 * 2. API处理：接收HTTP请求并调用服务层处理
 * 3. PocketBase集成：将服务作为PocketBase应用启动，并使用 OnServe 钩子挂载现有API
 *
 * 待优化事项：
 * 1. 配置管理：从环境变量或配置文件加载配置
 * 2. 错误处理：改进错误处理和恢复机制
 * 3. 日志记录：添加结构化日志记录
 *
 * 兼容性说明：
 * 1. 依赖service包中的服务接口和实现
 * 2. 兼容 PocketBase v0.23.0+
 */