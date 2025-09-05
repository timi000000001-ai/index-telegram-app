package main

import (
	"log"
	"net/http"

	"management-service/internal/config"
	_ "management-service/migrations"
	"management-service/service"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

	// Load configuration
	cfg, err := config.Load("configs")
	if err != nil {
		log.Fatal("Failed to load configuration: ", err)
	}

	// Initialize services
	searchSvc, botInfoSvc, webhookSvc := initServices(cfg)

	// Register API routes
	registerAPIs(app, searchSvc, botInfoSvc, webhookSvc)

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

func initServices(cfg *config.Config) (service.SearchService, service.BotInfoService, service.WebhookService) {
	searchConfig := service.SearchConfig{
		MeilisearchURL: cfg.MeilisearchURL,
		MeilisearchKey: cfg.MeilisearchKey,
	}
	searchService := service.NewSearchService(searchConfig)

	pocketBaseConfig := service.PocketBaseConfig{
		BaseURL: cfg.PocketBaseURL,
	}
	botInfoService := service.NewBotInfoService(pocketBaseConfig)

	webhookService := service.NewWebhookService(botInfoService, cfg.BotServiceURL)
	return searchService, botInfoService, webhookService
}

func registerAPIs(app *pocketbase.PocketBase, searchService service.SearchService, botInfoService service.BotInfoService, webhookService service.WebhookService) {
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		// Middleware to require admin authentication.
		// This function will be executed before each handler in the group.
		requireAdminAuth := func(e *core.RequestEvent) error {
			if !e.HasSuperuserAuth() {
				return apis.NewForbiddenError("Action requires admin access.", nil)
			}
			return nil // authorized
		}

		// Create a new router group for API endpoints and bind the middleware.
		apiGroup := se.Router.Group("/api")
		apiGroup.BindFunc(requireAdminAuth)

		// Register search API
		apiGroup.GET("/search", func(e *core.RequestEvent) error {
			q := e.Request.URL.Query().Get("q")
			p := e.Request.URL.Query().Get("page")
			l := e.Request.URL.Query().Get("limit")
			f := e.Request.URL.Query().Get("filter")
			results, err := searchService.Search(q, p, l, f)
			if err != nil {
				return apis.NewApiError(http.StatusInternalServerError, "Failed to perform search", err)
			}
			return e.JSON(http.StatusOK, results)
		})

		// Register bots API
		apiGroup.GET("/bots", func(e *core.RequestEvent) error {
			bots, err := botInfoService.GetAllBotInfos()
			//输出日志打印
			log.Printf("Successfully fetched %d bot(s)", len(bots))
			if err != nil {
				return apis.NewApiError(http.StatusInternalServerError, "Failed to get bot infos", err)
			}
			return e.JSON(http.StatusOK, bots)
		})

		// Register webhooks registration API
		apiGroup.POST("/webhooks/register", func(e *core.RequestEvent) error {
			if err := webhookService.RegisterWebhooks(); err != nil {
				return apis.NewApiError(http.StatusInternalServerError, "Failed to register webhooks", err)
			}
			return e.JSON(http.StatusOK, map[string]string{"status": "success"})
		})

		return se.Next()
	})
}
