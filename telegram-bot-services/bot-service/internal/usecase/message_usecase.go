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
}

type validationCacheEntry struct {
	isValid   bool
	timestamp time.Time
}

// messageUsecaseImpl 实现MessageUsecase接口
// @author fcj
// @date 2023-11-15
// @version 1.0.0
type messageUsecaseImpl struct {
	storageRepository repository.StorageRepository
	searchUsecase     SearchUsecase
	validationCache   map[string]validationCacheEntry
	cacheMutex        sync.RWMutex
}

// NewMessageUsecase 创建新的消息服务实例
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param storageRepository 存储仓库
// @param searchUsecase 搜索用例
// @return MessageUsecase 消息服务实例
func NewMessageUsecase(storageRepository repository.StorageRepository, searchUsecase SearchUsecase) MessageUsecase {
	return &messageUsecaseImpl{
		storageRepository: storageRepository,
		searchUsecase:     searchUsecase,
		validationCache:   make(map[string]validationCacheEntry),
	}
}

// SaveMessage 保存消息到存储系统
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param data 消息数据
// @return error 错误信息
func (m *messageUsecaseImpl) SaveMessage(data map[string]interface{}) error {
	// 使用存储服务的SaveAndIndex方法同时保存和索引消息
	if err := m.storageRepository.SaveAndIndex(data); err != nil {
		log.Printf("ERROR: failed to save and index message: %v", err)
		return fmt.Errorf("failed to save and index message: %w", err)
	}

	return nil
}

// BuildSearchResponse 构建搜索响应
func (m *messageUsecaseImpl) BuildSearchResponse(bot *telebot.Bot, query string, page int, filter string) (string, [][]telebot.InlineButton, error) {
	log.Printf("INFO: Building search response: query='%s', page=%d, filter='%s'", query, page, filter)
	if query == "" {
		return "", nil, errors.New("empty query")
	}
	// 打印开始搜索 query filter
	log.Printf("INFO: query: %s", query)
	log.Printf("INFO: filter: %s", filter)
	searchResult, err := m.searchUsecase.Search(query, page, 20, filter)
	//打印搜索完成
	if err != nil {
		log.Printf("ERROR: search failed: %v", err)
		return "", nil, fmt.Errorf("search failed: %w", err)
	}
	log.Printf("DEBUG: searchResult: %s", searchResult)

	var result struct {
		Hits       []map[string]interface{} `json:"hits"`
		Query      string                   `json:"query"`
		TotalPages int                      `json:"totalPages"`
		Page       int                      `json:"page"`
	}
	if err := json.Unmarshal(searchResult, &result); err != nil {
		log.Printf("ERROR: failed to unmarshal search result: %v", err)
		return "", nil, fmt.Errorf("failed to unmarshal search result: %w", err)
	}

	if len(result.Hits) == 0 {
		return "<i>No results found for: </i>" + html.EscapeString(query), nil, nil
	}

	// Temporarily disable live validation to avoid Telegram API rate limits.
	// The validation logic has been commented out.
	/*
		var validHits []map[string]interface{}
		for _, hit := range result.Hits {
			chatUsername, ok := hit["USERNAME"].(string)
			if ok && chatUsername != "" {
				m.cacheMutex.RLock()
				entry, exists := m.validationCache[chatUsername]
				m.cacheMutex.RUnlock()
	
				if exists && time.Since(entry.timestamp) < 24*time.Hour {
					if entry.isValid {
						validHits = append(validHits, hit)
					}
					continue
				}
				//此处需要处理 CCTAV1/16077 这种情况 只需要 CCTAV1 这种情况
				chatUsername = strings.Split(chatUsername, "/")[0]
				_, err := bot.ChatByUsername("@" + chatUsername)
				isValid := err == nil
				if err != nil {
					log.Printf("ERROR: Failed to get chat by username: %v", err)
				}
	
				m.cacheMutex.Lock()
				m.validationCache[chatUsername] = validationCacheEntry{
					isValid:   isValid,
					timestamp: time.Now(),
				}
				m.cacheMutex.Unlock()
	
				if isValid {
					validHits = append(validHits, hit)
				} else {
					log.Printf("INFO: Invalid username found and skipped: %s", chatUsername)
	
					// Send notification to review channel
	
					m.sendReviewNotification(bot, hit)
				}
			} else {
				// Keep results that don't have a username (e.g., private messages or chats without usernames)
				validHits = append(validHits, hit)
			}
		}
	
		if len(validHits) == 0 {
			return "<i>No valid results found for: </i>" + html.EscapeString(query), nil, nil
		}
		result.Hits = validHits
	*/

	response := fmt.Sprintf("<b>🔍 关键字: %s</b> (第 %d 页 / 共 %d 页)\n\n", html.EscapeString(query), result.Page, result.TotalPages)
	for i, hit := range result.Hits {
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
			// It's a message result
			messageID := int(messageIDFloat)
			messageText, textOk := hit["text"].(string)
			if textOk && messageText != "" {
				// Truncate for display
				if len([]rune(messageText)) > 120 {
					messageText = string([]rune(messageText)[:120]) + "..."
				}

				jumpLink := ""
				if chatUsername, ok := hit["USERNAME"].(string); ok && chatUsername != "" {
					jumpLink = fmt.Sprintf(" <a href=\"https://t.me/%s/%d\">(跳转)</a>", chatUsername, messageID)
				}

				response += fmt.Sprintf("<b>%d. 💬 消息</b> from %s%s\n", i+1+(page-1)*5, displayTitle, jumpLink)
				response += fmt.Sprintf("<blockquote>%s</blockquote>\n", html.EscapeString(messageText))
			}
		} else {
			// It's a chat/group/channel result
			var typeEmoji string
			if chatType, ok := hit["TYPE"].(string); ok {
				switch chatType {
				case "private":
					typeEmoji = "👤"
				//超级搜索显示更高端的符号
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

			response += fmt.Sprintf("<b>%d. %s</b> %s%s\n\n", i+1+(page-1)*5, displayTitle, typeEmoji, membersCountStr)
		}
	}
	var buttonRows [][]telebot.InlineButton

	// Row 1: Pagination
	paginationRow := []telebot.InlineButton{}
	if page > 1 {
		paginationRow = append(paginationRow, telebot.InlineButton{Text: "⬅️ 上一页", Data: fmt.Sprintf("prev_%s", filter)})
	}
	paginationRow = append(paginationRow, telebot.InlineButton{Text: fmt.Sprintf("%d/%d", result.Page, result.TotalPages), Data: "current"})
	if page < result.TotalPages {
		paginationRow = append(paginationRow, telebot.InlineButton{Text: "下一页 ➡️", Data: fmt.Sprintf("next_%s", filter)})
	}
	buttonRows = append(buttonRows, paginationRow)

	// Subsequent rows: Filters
	filterModels := []struct {
		Text  string
		Value string
	}{
		{"全部", "all"}, {"群组", "group"}, {"频道", "channel"}, {"机器人", "bot"}, {"消息", "message"},
	}

	var filterButtons []telebot.InlineButton
	for _, model := range filterModels {
		text := model.Text
		if (filter == "" && model.Value == "all") || filter == model.Value {
			text = "✅ " + text
		}
		filterButtons = append(filterButtons, telebot.InlineButton{Text: text, Data: fmt.Sprintf("filter_%s", model.Value)})
	}

	const maxButtonsPerRow = 3
	for i := 0; i < len(filterButtons); i += maxButtonsPerRow {
		end := i + maxButtonsPerRow
		if end > len(filterButtons) {
			end = len(filterButtons)
		}
		buttonRows = append(buttonRows, filterButtons[i:end])
	}

	log.Printf("Generated %d button rows for page %d, filter '%s'", len(buttonRows), page, filter)
	return response, buttonRows, nil
}

// SearchWithPagination 搜索消息并支持分页
func (m *messageUsecaseImpl) SearchWithPagination(c telebot.Context, query string, page int, filter string) error {
	bot, ok := c.Bot().(*telebot.Bot)
	if !ok {
		log.Printf("ERROR: could not get bot instance")
		return c.Send("内部错误，无法处理搜索请求。")
	}
	response, buttonRows, err := m.BuildSearchResponse(bot, query, page, filter)
	if err != nil {
		log.Printf("ERROR: Search failed: %v", err)
		return c.Send(fmt.Sprintf("🔍 搜索失败: `%s`", err.Error()), &telebot.SendOptions{ParseMode: telebot.ModeMarkdown})
	}
	if response == "" {
		return c.Send("请输入搜索关键字。")
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
		log.Printf("ERROR: could not get bot instance on callback")
		return c.Respond(&telebot.CallbackResponse{Text: "内部错误", ShowAlert: true})
	}

	newText, newMarkup, err := m.handleCallbackLogic(bot, c.Callback().Data, c.Callback().Message.Text)

	if err != nil {
		log.Printf("Error handling callback logic: %v", err)
		return c.Respond(&telebot.CallbackResponse{
			Text:      "操作失败: " + err.Error(),
			ShowAlert: false,
		})
	}

	// Case 1: Confirmation for delete/keep actions (newText is present, newMarkup is nil)
	if newText != "" && newMarkup == nil {
		// First, delete the original message after a delay.
		go func() {
			time.Sleep(2 * time.Second)
			if err := bot.Delete(c.Callback().Message); err != nil {
				log.Printf("ERROR: Failed to delete message on callback: %v", err)
			}
		}()

		// Then, show the alert. This also acknowledges the callback.
		return c.Respond(&telebot.CallbackResponse{
			Text:      newText,
			ShowAlert: false,
		})
	}

	// Case 2: Pagination or filter update (newText and newMarkup are present)
	if newText != "" {
		// c.Edit edits the message and acknowledges the callback in one go.
		return c.Edit(newText, &telebot.SendOptions{
			ParseMode:   telebot.ModeHTML,
			ReplyMarkup: &telebot.ReplyMarkup{InlineKeyboard: newMarkup},
		})
	}

	// Case 3: No-op (e.g., clicking current page)
	// newText and newMarkup are both empty.
	return c.Respond()
}

// handleCallbackLogic contains the testable logic for handling callbacks.
func (m *messageUsecaseImpl) handleCallbackLogic(bot *telebot.Bot, data, messageText string) (string, [][]telebot.InlineButton, error) {
	log.Printf("INFO: Handling callback logic: data='%s'", data)

	// First, handle review callbacks, as they don't have a query in their message.
	if strings.HasPrefix(data, "delete_doc_") {
		docID := strings.TrimPrefix(data, "delete_doc_")
		if err := m.searchUsecase.DeleteDocument(docID); err != nil {
			log.Printf("ERROR: failed to delete document %s: %v", docID, err)
			return "", nil, errors.New("删除失败")
		}
		// Return the confirmation message and nil buttons to remove the inline keyboard.
		return fmt.Sprintf("✅ 文档 %s 已被删除。", docID), nil, nil
	} else if strings.HasPrefix(data, "keep_doc_") {
		docID := strings.TrimPrefix(data, "keep_doc_")
		// Return the confirmation message and nil buttons to remove the inline keyboard.
		return fmt.Sprintf("👍 文档 %s 已被保留。", docID), nil, nil
	}

	var query, filter string
	var page, totalPages int

	// Extract query from the message text first, as it's always needed for search-related callbacks.
	reQuery := regexp.MustCompile(`<b>🔍 关键字: (.+?)</b>`)
	queryMatches := reQuery.FindStringSubmatch(messageText)
	if len(queryMatches) < 2 {
		// Fallback for older message formats without HTML
		re := regexp.MustCompile(`关键字: (\S+) \(`)
		fallbackMatches := re.FindStringSubmatch(messageText)
		if len(fallbackMatches) < 2 {
			log.Printf("ERROR: Could not extract query from message text: %s", messageText)
			return "", nil, errors.New("无法解析查询关键字")
		}
		query = fallbackMatches[1]
	} else {
		query = html.UnescapeString(queryMatches[1]) // Unescape HTML entities
	}

	// Handle different callback actions
	if strings.HasPrefix(data, "filter_") {
		filter = strings.TrimPrefix(data, "filter_")
		page = 1 // Reset to page 1 for new filter
	} else if strings.HasPrefix(data, "prev_") || strings.HasPrefix(data, "next_") {
		// Extract page and totalPages for pagination
		rePage := regexp.MustCompile(`\(第 (\d+) 页 / 共 (\d+) 页\)`)
		pageMatches := rePage.FindStringSubmatch(messageText)
		if len(pageMatches) < 3 {
			log.Printf("ERROR: Could not parse page number from message text: %s", messageText)
			return "", nil, errors.New("无法解析页码")
		}
		var err error
		page, err = strconv.Atoi(pageMatches[1])
		if err != nil {
			log.Printf("ERROR: Could not parse current page number: %v", err)
			return "", nil, errors.New("无法解析当前页码")
		}
		totalPages, err = strconv.Atoi(pageMatches[2])
		if err != nil {
			log.Printf("ERROR: Could not parse total pages number: %v", err)
			return "", nil, errors.New("无法解析总页码")
		}

		if strings.HasPrefix(data, "prev_") {
			filter = strings.TrimPrefix(data, "prev_")
			if page <= 1 {
				return "", nil, errors.New("已经是第一页了")
			}
			page--
		} else {
			filter = strings.TrimPrefix(data, "next_")
			if page >= totalPages {
				return "", nil, errors.New("已经是最后一页了")
			}
			page++
		}
	} else if data == "current" {
		return "", nil, nil // No-op, just acknowledge the callback
	} else {
		return "", nil, errors.New("未知操作")
	}

	// Build and send the updated message
	return m.BuildSearchResponse(bot, query, page, filter)
}

/*
 * 关键算法说明：
 * 1. 消息保存：将消息同时保存到PocketBase和索引到Meilisearch
 * 2. 分页搜索：支持按页码和过滤条件搜索消息
 * 3. 回调处理：处理用户与搜索结果交互的回调
 *
 * 待优化事项：
 * 1. 缓存机制：添加搜索结果缓存减少重复请求
 * 2. 批量处理：实现批量消息保存和索引
 *
 * 兼容性说明：
 * 1. 依赖telebot.v3库处理Telegram交互
 * 2. 需要配置正确的PocketBase和Meilisearch服务
 */

// sendReviewNotification sends a message to the review channel when an invalid username is found.
func (m *messageUsecaseImpl) sendReviewNotification(bot *telebot.Bot, hit map[string]interface{}) {
	chatUsername, ok := hit["USERNAME"].(string)
	if !ok {
		return // Should not happen if called correctly
	}

	// The document ID from MeiliSearch is in the 'id' field.
	idValue, exists := hit["id"]
	if !exists {
		log.Printf("ERROR: could not get document id (field 'id') for invalid username: %s", chatUsername)
		return
	}

	var docID string
	switch v := idValue.(type) {
	case string:
		docID = v
	case float64:
		// Telegram IDs are large numbers, convert float64 to int64 then to string
		docID = strconv.FormatInt(int64(v), 10)
	default:
		log.Printf("ERROR: 'id' has an unexpected type. Actual type: %T, value: %v", idValue, idValue)
		return
	}

	// First, resolve the username to a chat object to get the correct chat ID.
	reviewChat, err := bot.ChatByUsername("@SoSo00000000001")
	if err != nil {
		log.Printf("ERROR: could not find the review chat '@SoSo00000000001': %v", err)
		return
	}

	message := fmt.Sprintf(
		"***群组已经失效了已经删除***\n请审核是否真的无效: @%s",
		html.EscapeString(chatUsername),
	)

	inlineKeys := [][]telebot.InlineButton{
		{
			telebot.InlineButton{Text: "✅ 确认失效", Data: fmt.Sprintf("delete_doc_%s", docID)},
			telebot.InlineButton{Text: "❌ 有效", Data: fmt.Sprintf("keep_doc_%s", docID)},
		},
	}

	_, err = bot.Send(reviewChat, message, &telebot.SendOptions{
		ParseMode:   telebot.ModeHTML,
		ReplyMarkup: &telebot.ReplyMarkup{InlineKeyboard: inlineKeys},
	})
	if err != nil {
		log.Printf("ERROR: failed to send review notification for %s: %v", chatUsername, err)
	}
}
