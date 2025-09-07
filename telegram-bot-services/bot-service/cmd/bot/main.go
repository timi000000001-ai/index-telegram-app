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
	"strings"
	"sync"
	"time"

	"gopkg.in/telebot.v4"
)

var (
	tokenMutex     sync.Mutex
	tokenIndex     = 0
	tokenBlacklist = make(map[string]time.Time)
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

	messageUsecase := usecase.NewMessageUsecase(cfg, storageRepo, searchRepo)
	botHandler := handler.NewBotHandler(messageUsecase, cfg)

	// Initialize bots
	if err := initializeBots(botHandler, cfg); err != nil {
		log.Fatalf("Failed to initialize bots: %v", err)
	}

	// Initialize review bot
	if err := initializeReviewBot(botHandler, cfg); err != nil {
		log.Printf("Failed to initialize review bot: %v", err)
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
			if strings.Contains(err.Error(), "retry after") {
				tokenMutex.Lock()
				tokenBlacklist[botConfig.Token] = time.Now().Add(time.Duration(cfg.Bot.TokenRotationDuration) * time.Second)
				tokenMutex.Unlock()
			}
			continue
		}
		botHandler.RegisterHandlers(bot)
	}

	return nil
}

func initializeReviewBot(botHandler handler.BotHandler, cfg *config.Config) error {
	if cfg.Bot.ReviewBotToken == "" {
		return fmt.Errorf("review bot token is not configured")
	}

	bot, err := botHandler.InitBot(handler.BotConfig{
		Token: cfg.Bot.ReviewBotToken,
	}, cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize review bot: %w", err)
	}

	botHandler.RegisterReviewHandlers(bot)
	log.Println("Review bot initialized successfully")
	return nil
}

func getNextToken(cfg *config.Config) string {
	tokenMutex.Lock()
	defer tokenMutex.Unlock()

	for i := 0; i < len(cfg.Bot.BotTokens); i++ {
		token := cfg.Bot.BotTokens[tokenIndex]
		tokenIndex = (tokenIndex + 1) % len(cfg.Bot.BotTokens)

		if expiry, found := tokenBlacklist[token]; !found || time.Now().After(expiry) {
			return token
		}
	}

	return "" // No available token
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