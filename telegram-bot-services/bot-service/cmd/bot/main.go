package main

import (
	"bot-service/internal/api/handler"
	"bot-service/internal/config"
	"bot-service/internal/management"
	"bot-service/internal/repository"
	"bot-service/internal/usecase"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gopkg.in/telebot.v4"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("development")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize dependencies
	storageRepo := repository.NewStorageRepository(repository.StorageConfig{
		PocketBaseURL:    cfg.Storage.PocketBaseURL,
		MeilisearchURL:   cfg.Storage.MeilisearchURL,
		MeilisearchToken: cfg.Storage.MeilisearchToken,
	})

	searchRepo := repository.NewSearchRepository(cfg.Search.MeilisearchURL, cfg.Search.MeilisearchKey)
	searchUsecase := usecase.NewSearchUsecase(searchRepo)

	messageUsecase := usecase.NewMessageUsecase(storageRepo, searchUsecase)
	botHandler := handler.NewBotHandler(messageUsecase, cfg)

	// Initialize bots
	if err := initializeBots(botHandler, cfg); err != nil {
		log.Fatalf("Failed to initialize bots: %v", err)
	}

	// Start server
	startServer(botHandler)
}

func initializeBots(botHandler handler.BotHandler, cfg *config.Config) error {
	botConfigs, err := management.GetBotsToken()
	if err != nil {
		return fmt.Errorf("failed to get bot tokens: %w", err)
	}

	for _, botConfig := range botConfigs {
		bot, err := botHandler.InitBot(handler.BotConfig{
			Token:                  botConfig.Token,
			Name:                   botConfig.Name,
			Status:                 botConfig.Status,
			WebhookURL:             botConfig.WebhookURL,
			ManagementServiceURL:   botConfig.ManagementServiceURL,
			ManagementServiceToken: botConfig.ManagementServiceToken,
		}, cfg)
		if err != nil {
			log.Printf("Failed to initialize bot with token %s: %v", botConfig.Token, err)
			continue
		}
		botHandler.RegisterHandlers(bot)
	}

	return nil
}

func startServer(botHandler handler.BotHandler) {
	http.HandleFunc("/webhook", newWebhookHandler(botHandler))

	log.Println("Starting server on :8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func newWebhookHandler(botHandler handler.BotHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received webhook request: method=%s, url=%s", r.Method, r.URL.String())
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		token := r.URL.Query().Get("token")
		if token == "" {
			http.Error(w, "Token not found in URL query", http.StatusBadRequest)
			return
		}

		var update telebot.Update
		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			http.Error(w, "Failed to decode update", http.StatusBadRequest)
			return
		}

		if err := botHandler.ProcessUpdate(token, &update); err != nil {
			log.Printf("Failed to process update: %v", err)
			http.Error(w, "Failed to process update", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}