package main

import (
	"bot-service/internal/management"
	"bot-service/service"
	"encoding/json"
	"log"
	"net/http"

	"gopkg.in/telebot.v3"
)

func main() {
	if err := initializeBotService(); err != nil {
		log.Fatalf("Failed to initialize bot service: %v", err)
	}
	startServer()
}

func initializeBotService() error {
	botConfigs, err := management.GetBotsToken()
	if err != nil {
		return err
	}

	for _, botConfig := range botConfigs {
		if _, err := service.Init(botConfig); err != nil {
			log.Printf("Failed to initialize bot with token %s: %v", botConfig.Token, err)
			continue
		}
	}

	return nil
}

func startServer() {
	http.HandleFunc("/webhook", newWebhookHandler())

	log.Println("Starting server on :8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func newWebhookHandler() http.HandlerFunc {
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

		if err := service.ProcessUpdate(token, &update); err != nil {
			log.Printf("Failed to process update: %v", err)
			http.Error(w, "Failed to process update", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
