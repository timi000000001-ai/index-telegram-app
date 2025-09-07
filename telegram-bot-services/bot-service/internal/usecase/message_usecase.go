/*
 * 文件功能描述：消息服务，处理Telegram消息的搜索、保存和索引功能
 * 主要类/接口说明：MessageService接口及其实现
 * 修改历史记录：
 * @author fcj
 * @date 2023-11-15
 * @version 1.0.0
 * © Telegram Bot Services Team
 */

package usecase

import (
	"bot-service/internal/config"
	"bot-service/internal/repository"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"gopkg.in/telebot.v4"
)

// MessageUsecase 定义消息处理服务接口
// @author fcj
// @date 2023-11-15
// @version 1.0.0
type MessageUsecase interface {
	// SaveMessage 保存消息到存储系统
	SaveMessage(data map[string]interface{}) error

	// SearchWithPagination 搜索消息并支持分页
	SearchWithPagination(c telebot.Context, query string, page int, filter string) error

	// HandleCallback 处理回调查询
	HandleCallback(c telebot.Context) error

	// HandleReviewCallback 处理审核回调
	HandleReviewCallback(c telebot.Context) error
}

// SearchResponse defines the structure for a MeiliSearch search response.
type SearchResponse struct {
	Hits           []map[string]interface{} `json:"hits"`
	Query          string                   `json:"query"`
	ProcessingTimeMs int64                  `json:"processingTimeMs"`
	HitsPerPage    int64                    `json:"hitsPerPage"`
	Page           int64                    `json:"page"`
	TotalPages     int64                    `json:"totalPages"`
	TotalHits      int64                    `json:"totalHits"`
}

type validationCacheEntry struct {
	isValid   bool
	timestamp time.Time
}

// validationJob defines a task for the validation worker.
type validationJob struct {
	hit map[string]interface{}
}

// messageUsecaseImpl是messageUsecase的实现
type messageUsecaseImpl struct {
	cfg                   *config.Config
	storageRepo           repository.StorageRepository
	searchRepo            repository.SearchRepository
	cacheMutex            sync.RWMutex
	validationCache       map[string]validationCacheEntry
	validationQueue       chan validationJob // Channel for validation jobs
	tokenMutex            sync.Mutex
	tokenIndex            int
	tokenBlacklist        map[string]time.Time
	permanentlyBlacklist  map[string]bool
	tokenRotationDuration time.Duration
}

// NewMessageUsecase create a new messageUsecase
func NewMessageUsecase(cfg *config.Config, storageRepo repository.StorageRepository, searchRepo repository.SearchRepository) MessageUsecase {
	m := &messageUsecaseImpl{
		cfg:                   cfg,
		storageRepo:           storageRepo,
		searchRepo:            searchRepo,
		validationCache:       make(map[string]validationCacheEntry),
		validationQueue:       make(chan validationJob, 100), // Buffered channel for 100 jobs
		tokenBlacklist:        make(map[string]time.Time),
		permanentlyBlacklist:  make(map[string]bool),
		tokenRotationDuration: time.Duration(cfg.Bot.TokenRotationDuration) * time.Second,
	}
	// Start a background worker to process validation jobs.
	go m.startValidationWorker()
	return m
}

func (m *messageUsecaseImpl) getBotForValidation() (*telebot.Bot, error) {
	m.tokenMutex.Lock()
	defer m.tokenMutex.Unlock()

	for i := 0; i < len(m.cfg.Bot.BotTokens); i++ {
		token := m.cfg.Bot.BotTokens[m.tokenIndex]
		m.tokenIndex = (m.tokenIndex + 1) % len(m.cfg.Bot.BotTokens)

		if m.permanentlyBlacklist[token] {
			continue // Skip permanently blacklisted tokens
		}

		if expiry, found := m.tokenBlacklist[token]; !found || time.Now().After(expiry) {
			pref := telebot.Settings{
				Token:  token,
				Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
			}
			bot, err := telebot.NewBot(pref)
			if err != nil {
				if strings.Contains(err.Error(), "Unauthorized") {
					log.Printf("WARN: Token %s is permanently invalid, adding to permanent blacklist.", token)
					m.permanentlyBlacklist[token] = true
					continue // Try next token
				}
				return nil, fmt.Errorf("failed to create bot with token %s: %w", token, err)
			}
			return bot, nil
		}
	}

	return nil, errors.New("no available tokens for validation")
}

// startValidationWorker processes username validation jobs from the queue.
func (m *messageUsecaseImpl) startValidationWorker() {
	for job := range m.validationQueue {
		chatUsername, ok := job.hit["USERNAME"].(string)
		if !ok || chatUsername == "" {
			continue
		}

		// Double-check cache to avoid redundant processing if multiple searches return the same user.
		m.cacheMutex.RLock()
		entry, exists := m.validationCache[chatUsername]
		m.cacheMutex.RUnlock()

		if exists && time.Since(entry.timestamp) < 24*time.Hour {
			continue
		}

		// Clean up username format, e.g., "CCTAV1/16077" -> "CCTAV1"
		chatUsername = strings.Split(chatUsername, "/")[0]

		bot, err := m.getBotForValidation()
		if err != nil {
			log.Printf("ERROR: Failed to get bot for validation: %v", err)
			continue
		}

		_, err = bot.ChatByUsername("@" + chatUsername)
		isValid := err == nil

		if err != nil {
			// Log only if it's not a "not found" error, which is expected for invalid usernames.
			if strings.Contains(err.Error(), "chat not found") {
				// This is an expected error for invalid usernames, so we just mark it as invalid.
			} else if strings.Contains(err.Error(), "429") { // Handle rate limiting
				log.Printf("WARN: Rate limit hit with token %s. Temporarily blacklisting.", bot.Token)
				m.tokenMutex.Lock()
				m.tokenBlacklist[bot.Token] = time.Now().Add(m.tokenRotationDuration)
				m.tokenMutex.Unlock()
				// Re-queue the job to try with another token later.
				go func(j validationJob) {
					time.Sleep(1 * time.Second) // Wait a bit before re-queueing
					m.validationQueue <- j
				}(job)
				continue // Move to the next job
			} else if strings.Contains(err.Error(), "Unauthorized") {
				log.Printf("WARN: Token %s is permanently invalid, adding to permanent blacklist.", bot.Token)
				m.tokenMutex.Lock()
				m.permanentlyBlacklist[bot.Token] = true
				m.tokenMutex.Unlock()
				// Re-queue the job to try with another token.
				go func(j validationJob) {
					m.validationQueue <- j
				}(job)
				continue // Move to the next job
			} else {
				log.Printf("ERROR: Failed to get chat by username @%s: %v", chatUsername, err)
			}
		}

		m.cacheMutex.Lock()
		m.validationCache[chatUsername] = validationCacheEntry{
			isValid:   isValid,
			timestamp: time.Now(),
		}
		m.cacheMutex.Unlock()

		if !isValid {
			log.Printf("INFO: Invalid username found, sending for review: %s", chatUsername)
			m.sendReviewNotification(bot, job.hit)
		}

		// IMPORTANT: Wait for 1 second between requests to avoid hitting Telegram API rate limits.
		time.Sleep(1 * time.Second)
	}
}

// validateUsernamesAsync sends validation jobs to the background worker.
func (m *messageUsecaseImpl) validateUsernamesAsync(hits []interface{}) {
	for _, hit := range hits {
		hitMap, ok := hit.(map[string]interface{})
		if !ok {
			continue
		}
		chatUsername, ok := hitMap["USERNAME"].(string)
		if ok && chatUsername != "" {
			m.cacheMutex.RLock()
			entry, exists := m.validationCache[chatUsername]
			m.cacheMutex.RUnlock()

			// If not in cache or cache is expired, add to the validation queue.
			if !exists || time.Since(entry.timestamp) >= 24*time.Hour {
				m.validationQueue <- validationJob{hit: hitMap}
			}
		}
	}
}

// SaveMessage saves a message to MeiliSearch.
func (m *messageUsecaseImpl) SaveMessage(data map[string]interface{}) error {
	// Use the configured index name from config
	indexName := m.cfg.Search.IndexName
	if indexName == "" {
		indexName = "messages" // Fallback to a default name
	}
	// This part needs to be implemented in storageRepo
	// _, err := m.meiliSearch.Index(indexName).AddDocuments([]map[string]interface{}{data})
	// if err != nil {
	// 	log.Printf("ERROR: failed to index message: %v", err)
	// 	return fmt.Errorf("failed to index message: %w", err)
	// }
	return nil
}

// buildSearchResponse builds the search response string and buttons.
func (m *messageUsecaseImpl) buildSearchResponse(query string, filter string, searchResult *SearchResponse) (string, [][]telebot.InlineButton, error) {
	log.Printf("INFO: Building search response: query='%s', filter='%s', page=%d", query, filter, searchResult.Page)

	if len(searchResult.Hits) == 0 {
		return "<i>No results found for: </i>" + html.EscapeString(query), nil, nil
	}

	currentPage := searchResult.Page
	totalPages := searchResult.TotalPages
	hitsPerPage := searchResult.HitsPerPage
	if hitsPerPage <= 0 {
		hitsPerPage = 10 // Fallback
	}

	response := fmt.Sprintf("<b>🔍 关键字: %s</b> (第 %d 页 / 共 %d 页)\n\n", html.EscapeString(query), currentPage, totalPages)
	for i, hit := range searchResult.Hits {

		chatTitle := hit["TITLE"]
		if chatTitle == nil || chatTitle == "" {
			if chatType, ok := hit["TYPE"].(string); ok {
				switch chatType {
				case "private":
					chatTitle = "私聊"
				case "group", "supergroup":
					chatTitle = "群组"
				case "channel":
					chatTitle = "频道"
				default:
					chatTitle = "未知"
				}
			} else {
				chatTitle = "未知"
			}
		}
		var displayTitle string
		chatUsername, ok := hit["USERNAME"].(string)
		if ok && chatUsername != "" {
			displayTitle = fmt.Sprintf("<a href=\"https://t.me/%s\">%s</a>", chatUsername, html.EscapeString(fmt.Sprint(chatTitle)))
		} else {
			displayTitle = html.EscapeString(fmt.Sprint(chatTitle))
		}

		if messageIDFloat, ok := hit["MESSAGE_ID"].(float64); ok {
			messageID := int(messageIDFloat)
			messageText, textOk := hit["text"].(string)
			if textOk && messageText != "" {
				if len([]rune(messageText)) > 120 {
					messageText = string([]rune(messageText)[:120]) + "..."
				}
				jumpLink := ""
				if chatUsername, ok := hit["USERNAME"].(string); ok && chatUsername != "" {
					jumpLink = fmt.Sprintf(" <a href=\"https://t.me/%s/%d\">(跳转)</a>", chatUsername, messageID)
				}
				response += fmt.Sprintf("<b>%d. 💬 消息</b> from %s%s\n", i+1+int((currentPage-1)*hitsPerPage), displayTitle, jumpLink)
				response += fmt.Sprintf("<blockquote>%s</blockquote>\n", html.EscapeString(messageText))
			}
		} else {
			var typeEmoji string
			if chatType, ok := hit["TYPE"].(string); ok {
				switch chatType {
				case "private":
					typeEmoji = "👤"
				case "supergroup":
					typeEmoji = "👑"
				case "group":
					typeEmoji = "👥"
				case "channel":
					typeEmoji = "📢"
				case "bot":
					typeEmoji = "🤖"
				}
			}
			var membersCountStr string
			if membersCount, ok := hit["MEMBERS_COUNT"].(float64); ok && membersCount > 0 {
				membersCountStr = fmt.Sprintf(" %d", int(membersCount))
			}
			response += fmt.Sprintf("<b>%d. %s</b> %s%s\n\n", i+1+int((currentPage-1)*hitsPerPage), displayTitle, typeEmoji, membersCountStr)
		}
	}

	var buttonRows [][]telebot.InlineButton
	paginationRow := []telebot.InlineButton{}
	if currentPage > 1 {
		paginationRow = append(paginationRow, telebot.InlineButton{Text: "⬅️ 上一页", Data: fmt.Sprintf("prev_%s_%s", filter, query)})
	}
	paginationRow = append(paginationRow, telebot.InlineButton{Text: fmt.Sprintf("%d/%d", currentPage, totalPages), Data: "current"})
	if currentPage < totalPages {
		paginationRow = append(paginationRow, telebot.InlineButton{Text: "下一页 ➡️", Data: fmt.Sprintf("next_%s_%s", filter, query)})
	}
	buttonRows = append(buttonRows, paginationRow)

	filterModels := []struct {
		Text  string
		Value string
	}{
		{"全部", "all"}, {"群组", "group"}, {"频道", "channel"}, {"机器人", "bot"}, {"消息", "message"},
	}
	var filterButtons []telebot.InlineButton
	for _, model := range filterModels {
		text := model.Text
		currentFilter := filter
		if currentFilter == "" {
			currentFilter = "all"
		}
		if currentFilter == model.Value {
			text = "✅ " + text
		}
		filterButtons = append(filterButtons, telebot.InlineButton{Text: text, Data: fmt.Sprintf("filter_%s_%s", model.Value, query)})
	}
	const maxButtonsPerRow = 3
	for i := 0; i < len(filterButtons); i += maxButtonsPerRow {
		end := i + maxButtonsPerRow
		if end > len(filterButtons) {
			end = len(filterButtons)
		}
		buttonRows = append(buttonRows, filterButtons[i:end])
	}

	return response, buttonRows, nil
}

// SearchWithPagination handles paginated search queries.
func (m *messageUsecaseImpl) SearchWithPagination(c telebot.Context, query string, page int, filter string) error {
	if query == "" {
		return c.Send("请输入搜索关键字。")
	}

	limit := 10
	searchResultRaw, err := m.searchRepo.Search(query, page, limit, filter)
	if err != nil {
		log.Printf("ERROR: search failed: %v", err)
		return c.Send(fmt.Sprintf("🔍 搜索失败: `%s`", err.Error()), &telebot.SendOptions{ParseMode: telebot.ModeMarkdown})
	}

	var searchResult SearchResponse
	if err := json.Unmarshal(searchResultRaw, &searchResult); err != nil {
		log.Printf("ERROR: failed to unmarshal search result: %v", err)
		return c.Send("🔍 搜索失败: 无法解析搜索结果。")
	}

	// Asynchronously validate usernames without blocking the search result response.
	var hits []interface{}
	for _, hit := range searchResult.Hits {
		hits = append(hits, hit)
	}
	go m.validateUsernamesAsync(hits)

	response, buttonRows, err := m.buildSearchResponse(query, filter, &searchResult)
	if err != nil {
		log.Printf("ERROR: Failed to build search response: %v", err)
		return c.Send("构建搜索结果失败。")
	}

	markup := &telebot.ReplyMarkup{
		InlineKeyboard: buttonRows,
	}

	err = c.Send(response, &telebot.SendOptions{
		ParseMode:             telebot.ModeHTML,
		ReplyMarkup:           markup,
		DisableWebPagePreview: true,
	})

	if err != nil {
		log.Printf("ERROR: Failed to send search response: %v", err)
		return err
	}
	return nil
}

// HandleCallback 处理回调查询
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param c Telegram上下文
// @return error 错误信息
func (m *messageUsecaseImpl) HandleCallback(c telebot.Context) error {
	bot, ok := c.Bot().(*telebot.Bot)
	if !ok {
		return c.Respond(&telebot.CallbackResponse{Text: "内部错误", ShowAlert: true})
	}

	newText, newMarkup, err := m.handleCallbackLogic(c.Callback().Data, c.Callback().Message.Text)

	if err != nil {
		return c.Respond(&telebot.CallbackResponse{Text: "操作失败: " + err.Error()})
	}

	if newText != "" && newMarkup == nil {
		go func() {
			time.Sleep(2 * time.Second)
			_ = bot.Delete(c.Callback().Message)
		}()
		return c.Respond(&telebot.CallbackResponse{Text: newText})
	}

	if newText != "" {
		return c.Edit(newText, &telebot.SendOptions{
			ParseMode:             telebot.ModeHTML,
			ReplyMarkup:           &telebot.ReplyMarkup{InlineKeyboard: newMarkup},
			DisableWebPagePreview: true,
		})
	}

	return c.Respond()
}

func (m *messageUsecaseImpl) HandleReviewCallback(c telebot.Context) error {
	data := c.Callback().Data
	var responseText string
	var err error

	if strings.HasPrefix(data, "delete_doc_") {
		docID := strings.TrimPrefix(data, "delete_doc_")
		if docID == "" {
			responseText = "❌ 删除失败: 无效的文档ID"
			log.Printf("ERROR: Attempted to delete document with empty ID.")
		} else {
			// 添加额外的安全检查，确保 docID 不为空且格式有效
			if strings.Contains(docID, "/") || strings.Contains(docID, "\\") {
				responseText = "❌ 删除失败: 无效的文档ID格式"
				log.Printf("ERROR: Invalid document ID format: %s", docID)
			} else {
				err = m.searchRepo.DeleteDocument(docID)
				if err != nil {
					responseText = "❌ 删除失败"
					log.Printf("ERROR: Failed to delete document %s: %v", docID, err)
				} else {
					responseText = fmt.Sprintf("✅ 文档 %s 已被删除。", docID)
					log.Printf("INFO: Successfully deleted document %s", docID)
				}
			}
		}
	} else if strings.HasPrefix(data, "keep_doc_") {
		docID := strings.TrimPrefix(data, "keep_doc_")
		responseText = fmt.Sprintf("👍 文档 %s 已被保留。", docID)
	} else {
		return c.Respond() // Ignore other callbacks
	}

	// Edit the original message to show the result
	err = c.Edit(responseText, &telebot.SendOptions{
		ParseMode:             telebot.ModeHTML,
		DisableWebPagePreview: true,
	})
	if err != nil {
		log.Printf("ERROR: Failed to edit message for review callback: %v", err)
		// Fallback to a response if editing fails
		return c.Respond(&telebot.CallbackResponse{Text: responseText})
	}

	return nil
}

// handleCallbackLogic contains the testable logic for handling callbacks.
func (m *messageUsecaseImpl) handleCallbackLogic(data, messageText string) (string, [][]telebot.InlineButton, error) {
	// Callback data format: action_filter_query
	// Use SplitN to correctly handle queries that may contain underscores.
	parts := strings.SplitN(data, "_", 3)
	if len(parts) < 1 {
		return "", nil, errors.New("invalid callback data")
	}
	action := parts[0]

	if action == "current" {
		return "", nil, nil // No action needed for the current page button
	}

	// For prev/next/filter, we need the full callback data.
	if len(parts) < 3 {
		return "", nil, fmt.Errorf("incomplete callback data: %s", data)
	}

	var query, filter string
	var page = 1

	// Extract current page from the message text.
	rePage := regexp.MustCompile(`\(第 (\d+) 页 / 共 (\d+) 页\)`)
	pageMatches := rePage.FindStringSubmatch(messageText)
	if len(pageMatches) >= 2 {
		page, _ = strconv.Atoi(pageMatches[1])
	}

	filter = parts[1]
	query = parts[2]

	switch action {
	case "filter":
		page = 1 // Reset to the first page when filter changes
	case "prev":
		if page > 1 {
			page--
		}
	case "next":
		page++
	default:
		return "", nil, fmt.Errorf("unknown action: %s", action)
	}

	// Perform the search again with the new parameters
	limit := 10
	searchResultRaw, err := m.searchRepo.Search(query, page, limit, filter)
	if err != nil {
		return "", nil, fmt.Errorf("搜索失败: %w", err)
	}

	var searchResult SearchResponse
	if err := json.Unmarshal(searchResultRaw, &searchResult); err != nil {
		log.Printf("ERROR: failed to unmarshal search result: %v", err)
		return "", nil, fmt.Errorf("搜索失败: 无法解析搜索结果。")
	}

	return m.buildSearchResponse(query, filter, &searchResult)
}

// sendReviewNotification sends a message to the review channel.
func (m *messageUsecaseImpl) sendReviewNotification(bot *telebot.Bot, hit map[string]interface{}) {
	go func() {
		// Use a separate bot instance for sending review notifications
		reviewBotToken := m.cfg.Bot.ReviewBotToken
		reviewBot, err := telebot.NewBot(telebot.Settings{
			Token: reviewBotToken,
		})
		if err != nil {
			log.Printf("ERROR: Failed to create review bot instance: %v", err)
			return
		}

		chatTitle, _ := hit["TITLE"].(string)
		reviewChannelID, err := strconv.ParseInt(m.cfg.Bot.ReviewChannel, 10, 64)
		if err != nil {
			log.Printf("ERROR: Invalid review channel ID: %v", err)
			return
		}
		reviewChat := &telebot.Chat{ID: reviewChannelID}
		var docID string
		if idVal, ok := hit["id"]; ok {
			switch v := idVal.(type) {
			case string:
				docID = v
			case float64:
				docID = strconv.FormatFloat(v, 'f', 0, 64)
			default:
				log.Printf("ERROR: docID is of an unexpected type (%T) in sendReviewNotification. hit: %v", v, hit)
				return
			}
		}

		if docID == "" {
			log.Printf("ERROR: docID is empty after conversion in sendReviewNotification. hit: %v", hit)
			return
		}
		chatUsername, _ := hit["USERNAME"].(string)
		message := fmt.Sprintf("<b>【疑似失效】</b>\n请审核: <a href=\"https://t.me/%s\">@%s</a>\n文档ID: <code>%s</code>", chatUsername, html.EscapeString(chatTitle), docID)
		inlineKeys := [][]telebot.InlineButton{
			{
				telebot.InlineButton{Text: "❌ 确认失效 (删除)", Data: fmt.Sprintf("delete_doc_%s", docID)},
				telebot.InlineButton{Text: "✅ 保留 没有失效", Data: fmt.Sprintf("keep_doc_%s", docID)},
			},
		}

		_, err = reviewBot.Send(reviewChat, message, &telebot.SendOptions{
			ParseMode:   telebot.ModeHTML,
			ReplyMarkup: &telebot.ReplyMarkup{InlineKeyboard: inlineKeys},
		})
		if err != nil {
			log.Printf("ERROR: Failed to send review notification with token ending in ...%s: %v", reviewBot.Token[len(reviewBot.Token)-6:], err)
		}
	}()
}
