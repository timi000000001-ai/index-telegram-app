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
    "fmt"
    "strings"
    "time"

    "gopkg.in/telebot.v3"
)

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

// SearchWithPagination 搜索消息并支持分页
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param c Telegram上下文
// @param query 搜索查询
// @param page 页码
// @param filter 过滤条件
// @return error 错误信息
func (m *messageServiceImpl) SearchWithPagination(c telebot.Context, query string, page int, filter string) error {
    // 调用搜索服务进行搜索
    searchResult, err := m.searchService.Search(query, page, 5, filter)
    if err != nil {
        return c.Send("Search failed: " + err.Error())
    }
    
    // 解析搜索结果
    var result struct {
        Hits      []map[string]interface{} `json:"hits"`
        TotalHits int                      `json:"totalHits"`
        TotalPages int                      `json:"totalPages"`
    }
    json.Unmarshal(searchResult, &result)
    
    if len(result.Hits) == 0 {
        return c.Send("No results found for: " + query)
    }

    // 构建响应消息
    response := fmt.Sprintf("Results for '%s' (Page %d of %d) - %s:\n", query, page, result.TotalPages, time.Now().Format("2006-01-02 15:04:05 +07"))
    for _, hit := range result.Hits {
        response += fmt.Sprintf("%s: %s (%s)\n", hit["chat_title"], hit["text"], hit["date"])
    }

    // 构建分页按钮
    var buttons []telebot.InlineButton
    if page > 1 {
        buttons = append(buttons, telebot.InlineButton{Text: "🔙 Previous", Data: fmt.Sprintf("prev_%d_%s", page, filter)})
    }
    if page < result.TotalPages {
        buttons = append(buttons, telebot.InlineButton{Text: "🔜 Next", Data: fmt.Sprintf("next_%d_%s", page, filter)})
    }

    // 构建过滤按钮
    response += "\nFilter by: "
    filterButtons := []struct {
        Text   string
        Value  string
        Active bool
    }{
        {"📢 Group", "group", filter == "group"},
        {"📡 Channel", "channel", filter == "channel"},
        {"🤖 Bot", "bot", filter == "bot"},
        {"💬 Message", "message", filter == "message"},
    }
    for _, btn := range filterButtons {
        text := btn.Text
        if btn.Active {
            text = "✅ " + text
        }
        buttons = append(buttons, telebot.InlineButton{Text: text, Data: fmt.Sprintf("filter_%s_%d", btn.Value, page)})
    }

    // 将按钮分组，每行最多显示2个按钮，优化UI布局
    var buttonRows [][]telebot.InlineButton
    const buttonsPerRow = 2
    for i := 0; i < len(buttons); i += buttonsPerRow {
        end := i + buttonsPerRow
        if end > len(buttons) {
            end = len(buttons)
        }
        buttonRows = append(buttonRows, buttons[i:end])
    }
    
    // 创建正确的ReplyMarkup对象
    markup := &telebot.ReplyMarkup{}
    markup.InlineKeyboard = buttonRows
    
    return c.Send(response, markup)
}

// HandleCallback 处理回调查询
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param c Telegram上下文
// @return error 错误信息
func (m *messageServiceImpl) HandleCallback(c telebot.Context) error {
    data := c.Callback().Data
    parts := strings.Split(data, "_")
    if len(parts) < 2 {
        return c.Respond(&telebot.CallbackResponse{Text: "Invalid action"})
    }

    action, value := parts[0], parts[1]
    var page int
    var filter string
    fmt.Sscanf(value, "%d_%s", &page, &filter)
    if filter == "" {
        filter = "message"
    }

    switch action {
    case "prev":
        page--
    case "next":
        page++
    case "filter":
        filter = parts[1]
    }

    query := c.Message().Text
    return m.SearchWithPagination(c, query, page, filter)
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