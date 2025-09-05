package service

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

// WebhookService defines the interface for managing Telegram bot webhooks.
// @author fcj
// @date 2023-11-16
// @version 1.0.0
type WebhookService interface {
	RegisterWebhooks() error
}

// webhookServiceImpl implements the WebhookService interface.
// @author fcj
// @date 2023-11-16
// @version 1.0.0
type webhookServiceImpl struct {
	client         *resty.Client
	botInfoService BotInfoService
	botServiceURL  string
}

// NewWebhookService creates a new WebhookService instance.
// @author fcj
// @date 2023-11-16
// @version 1.0.0
// @param botInfoService The service for fetching bot information.
// @param botServiceURL The base URL of the bot service.
// @return WebhookService A new WebhookService instance.
func NewWebhookService(botInfoService BotInfoService, botServiceURL string) WebhookService {
	return &webhookServiceImpl{
		client:         resty.New(),
		botInfoService: botInfoService,
		botServiceURL:  botServiceURL,
	}
}

// RegisterWebhooks fetches all bot tokens and registers their webhooks with Telegram.
// @author fcj
// @date 2023-11-16
// @version 1.0.0
// @return error An error if the registration process fails.
func (s *webhookServiceImpl) RegisterWebhooks() error {
	bots, err := s.botInfoService.GetAllBotInfos()
	if err != nil {
		return fmt.Errorf("failed to get bot infos: %w", err)
	}

	for _, bot := range bots {
		webhookURL := fmt.Sprintf("%s/webhook/%s", s.botServiceURL, bot.Token)
		telegramAPIURL := fmt.Sprintf("http：127.0.0.1:8082/bot%s/setWebhook", bot.Token)

		resp, err := s.client.R().
			SetBody(map[string]string{"url": webhookURL}).
			Post(telegramAPIURL)

		if err != nil {
			return fmt.Errorf("failed to set webhook for bot %s: %w", bot.ID, err)
		}

		if resp.IsError() {
			return fmt.Errorf("failed to set webhook for bot %s: status code %d, body: %s", bot.ID, resp.StatusCode(), resp.String())
		}
	}

	return nil
}