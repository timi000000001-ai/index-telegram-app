/*
 * 文件功能描述：消息服务，处理Telegram消息的搜索、保存和索引功能
 * 主要类/接口说明：MessageService接口及其实现
 * 修改历史记录：
 * @author fcj
 * @date 2023-11-15
 * @version 1.0.0
 * © Telegram Bot Services Team
 */

package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"gopkg.in/telebot.v3"
)

func escapeHTML(text string) string {
	text = strings.ReplaceAll(text, "&", "&amp;")
	text = strings.ReplaceAll(text, "<", "&lt;")
	text = strings.ReplaceAll(text, ">", "&gt;")
	text = strings.ReplaceAll(text, "\"", "&quot;")
	return text
}

// MessageService 定义消息处理服务接口
// @author fcj
// @date 2023-11-15
// @version 1.0.0
type MessageService interface {
	// SaveMessage 保存消息到存储系统
	SaveMessage(data map[string]interface{}) error

	// SearchWithPagination 搜索消息并支持分页
	SearchWithPagination(c telebot.Context, query string, page int, filter string) error

	// HandleCallback 处理回调查询
	HandleCallback(c telebot.Context) error
}

// messageServiceImpl 实现MessageService接口
// @author fcj
// @date 2023-11-15
// @version 1.0.0
type messageServiceImpl struct {
	storageService StorageService
	searchService  SearchService
}

// NewMessageService 创建新的消息服务实例
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param storageService 存储服务
// @param searchService 搜索服务
// @return MessageService 消息服务实例
func NewMessageService(storageService StorageService, searchService SearchService) MessageService {
	return &messageServiceImpl{
		storageService: storageService,
		searchService:  searchService,
	}
}

// SaveMessage 保存消息到存储系统
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param data 消息数据
// @return error 错误信息
func (m *messageServiceImpl) SaveMessage(data map[string]interface{}) error {
	// 使用存储服务的SaveAndIndex方法同时保存和索引消息
	if err := m.storageService.SaveAndIndex(data); err != nil {
		return fmt.Errorf("failed to save and index message: %w", err)
	}

	return nil
}

// BuildSearchResponse 构建搜索响应
func (m *messageServiceImpl) BuildSearchResponse(query string, page int, filter string) (string, [][]telebot.InlineButton, error) {
	log.Printf("Building search response: query='%s', page=%d, filter='%s'", query, page, filter)
	if query == "" {
		return "", nil, errors.New("empty query")
	}
	searchResult, err := m.searchService.Search(query, page, 5, filter)
	if err != nil {
		return "", nil, err
	}
	var result struct {
		Hits               []map[string]interface{} `json:"hits"`
		EstimatedTotalHits int                      `json:"estimatedTotalHits"`
		Limit              int                      `json:"limit"`
	}
	if err := json.Unmarshal(searchResult, &result); err != nil {
		return "", nil, err
	}
	totalPages := 0
	if result.Limit > 0 {
		totalPages = (result.EstimatedTotalHits + result.Limit - 1) / result.Limit
	}
	if len(result.Hits) == 0 {
		return "No results found for: " + query, nil, nil
	}
	response := fmt.Sprintf("<b>🔍 关键字: %s</b> (第 %d 页 / 共 %d 页)\n\n", escapeHTML(query), page, totalPages)
	for i, hit := range result.Hits {
		chatTitle := hit["title"]
		if chatTitle == nil || chatTitle == "" {
			if chatType, ok := hit["chat_type"].(string); ok {
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
		username, ok := hit["username"].(string)
		if ok && username != "" {
			displayTitle = fmt.Sprintf("<a href=\"https://t.me/%s\">%s</a>", username, escapeHTML(fmt.Sprint(chatTitle)))
		} else {
			displayTitle = escapeHTML(fmt.Sprint(chatTitle))
		}
		text := "无消息文本"
		if hit["text"] != nil {
			text = fmt.Sprint(hit["text"])
		}
		date := "未知日期"
		if hit["date"] != nil {
			date = fmt.Sprint(hit["date"])
		}
		response += fmt.Sprintf("<b>%d. %s</b>\n<i>%s</i>\n📅 %s\n\n", i+1+(page-1)*5, displayTitle, escapeHTML(text), escapeHTML(date))
	}
	var buttons []telebot.InlineButton
	if page > 1 {
		buttons = append(buttons, telebot.InlineButton{Text: "上一页", Data: fmt.Sprintf("prev_%d_%s", page, filter)})
	}
	buttons = append(buttons, telebot.InlineButton{Text: fmt.Sprintf("%d/%d", page, totalPages), Data: "current"})
	if page < totalPages {
		buttons = append(buttons, telebot.InlineButton{Text: "下一页", Data: fmt.Sprintf("next_%d_%s", page, filter)})
	}
	filterButtons := []struct {
		Text   string
		Value  string
		Active bool
	}{
		{"全部", "all", filter == "" || filter == "all"},
		{"群组", "group", filter == "group"},
		{"频道", "channel", filter == "channel"},
		{"机器人", "bot", filter == "bot"},
		{"消息", "message", filter == "message"},
	}
	for _, btn := range filterButtons {
		text := btn.Text
		if btn.Active {
			text = "✅ " + text
		}
		buttons = append(buttons, telebot.InlineButton{Text: text, Data: fmt.Sprintf("filter_%s_%d", btn.Value, page)})
	}
	var buttonRows [][]telebot.InlineButton
	buttonRows = append(buttonRows, buttons[:3])
	for i := 3; i < len(buttons); i += 3 {
		end := i + 3
		if end > len(buttons) {
			end = len(buttons)
		}
		buttonRows = append(buttonRows, buttons[i:end])
	}
	log.Printf("Generated %d button rows for page %d, filter '%s'", len(buttonRows), page, filter)
	return response, buttonRows, nil
}

// SearchWithPagination 搜索消息并支持分页
func (m *messageServiceImpl) SearchWithPagination(c telebot.Context, query string, page int, filter string) error {
	response, buttonRows, err := m.BuildSearchResponse(query, page, filter)
	if err != nil {
		return c.Send("Search failed: " + err.Error())
	}
	if response == "" {
		return c.Send("Please provide a search query.")
	}
	//打印下markup
	markup := &telebot.ReplyMarkup{}
	markup.InlineKeyboard = buttonRows
	err = c.Send(response, markup, &telebot.SendOptions{ParseMode: telebot.ModeHTML})
	if err != nil {
		log.Printf("Failed to send search response: %v", err)
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
func (m *messageServiceImpl) HandleCallback(c telebot.Context) error {
	log.Printf("Handling callback: data='%s'", c.Callback().Data)
	data := c.Callback().Data
	parts := strings.Split(data, "_")
	if len(parts) < 1 {
		return c.Respond(&telebot.CallbackResponse{Text: "Invalid action"})
	}

	action := parts[0]
	var page int
	var filter string
	var err error

	switch action {
	case "prev", "next":
		if len(parts) != 3 {
			return c.Respond(&telebot.CallbackResponse{Text: "Invalid data"})
		}
		page, err = strconv.Atoi(parts[1])
		if err != nil {
			return c.Respond(&telebot.CallbackResponse{Text: "Invalid page"})
		}
		filter = parts[2]
		if action == "prev" {
			page--
		} else {
			page++
		}
	case "filter":
		if len(parts) != 3 {
			return c.Respond(&telebot.CallbackResponse{Text: "Invalid data"})
		}
		filter = parts[1]
		page, err = strconv.Atoi(parts[2])
		if err != nil {
			return c.Respond(&telebot.CallbackResponse{Text: "Invalid page"})
		}
	case "current":
		return c.Respond(&telebot.CallbackResponse{Text: "当前页"})
	default:
		return c.Respond(&telebot.CallbackResponse{Text: "Unknown action"})
	}

	if page < 1 {
		page = 1
	}
	if filter == "" {
		filter = "message"
	}

	query := strings.TrimPrefix(c.Message().Text, "/search ")
	if query == "" {
		query = c.Message().Text
	}
	response, markup, err := m.BuildSearchResponse(query, page, filter)
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{Text: "Search failed"})
	}
	return c.Edit(response, markup, &telebot.SendOptions{ParseMode: telebot.ModeHTML})
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
