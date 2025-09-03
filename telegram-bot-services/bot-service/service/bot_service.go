/*
 * 文件功能描述：机器人服务，处理Telegram机器人的初始化、配置和管理
 * 主要类/接口说明：BotService接口及其实现
 * 修改历史记录：
 * @author fcj
 * @date 2023-11-15
 * @version 1.0.0
 * © Telegram Bot Services Team
 */

package service

import (
    "fmt"
    "sync"

    "gopkg.in/telebot.v3"
)

// BotConfig 定义机器人配置
// @author fcj
// @date 2023-11-15
// @version 1.0.0
type BotConfig struct {
    Token      string `json:"token"`
    WebhookURL string `json:"webhook_url"`
}

// BotService 定义机器人服务接口
// @author fcj
// @date 2023-11-15
// @version 1.0.0
type BotService interface {
    // InitBot 初始化机器人
    InitBot(config BotConfig) error
    
    // GetBot 获取机器人实例
    GetBot(token string) (*telebot.Bot, bool)
    
    // ProcessUpdate 处理Webhook更新
    ProcessUpdate(token string, update *telebot.Update) error
    
    // RegisterHandlers 注册消息处理函数
    RegisterHandlers(bot *telebot.Bot, messageService MessageService)
}

// botServiceImpl 实现BotService接口
// @author fcj
// @date 2023-11-15
// @version 1.0.0
type botServiceImpl struct {
    bots    map[string]*telebot.Bot
    mutex   sync.RWMutex
    configs []BotConfig
}

// NewBotService 创建新的机器人服务实例
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @return BotService 机器人服务实例
func NewBotService() BotService {
    return &botServiceImpl{
        bots:  make(map[string]*telebot.Bot),
        mutex: sync.RWMutex{},
    }
}

// InitBot 初始化机器人
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param config 机器人配置
// @return error 错误信息
func (b *botServiceImpl) InitBot(config BotConfig) error {
    bot, err := telebot.NewBot(telebot.Settings{
        Token: config.Token,
        Poller: &telebot.Webhook{
            Listen:   ":8081",
            Endpoint: &telebot.WebhookEndpoint{PublicURL: config.WebhookURL},
        },
    })
    if err != nil {
        return fmt.Errorf("failed to init bot %s: %v", config.Token, err)
    }

    webhook := &telebot.Webhook{Endpoint: &telebot.WebhookEndpoint{PublicURL: config.WebhookURL}}
    if err := bot.SetWebhook(webhook); err != nil {
        return fmt.Errorf("failed to set webhook for %s: %v", config.Token, err)
    }

    b.mutex.Lock()
    b.bots[config.Token] = bot
    b.configs = append(b.configs, config)
    b.mutex.Unlock()

    go bot.Start()
    return nil
}

// GetBot 获取机器人实例
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param token 机器人令牌
// @return *telebot.Bot 机器人实例
// @return bool 是否存在
func (b *botServiceImpl) GetBot(token string) (*telebot.Bot, bool) {
    b.mutex.RLock()
    defer b.mutex.RUnlock()
    
    bot, exists := b.bots[token]
    return bot, exists
}

// ProcessUpdate 处理Webhook更新
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param token 机器人令牌
// @param update Webhook更新
// @return error 错误信息
func (b *botServiceImpl) ProcessUpdate(token string, update *telebot.Update) error {
    bot, exists := b.GetBot(token)
    if !exists {
        return fmt.Errorf("bot not found: %s", token)
    }
    
    bot.ProcessUpdate(*update)
    return nil
}

// RegisterHandlers 注册消息处理函数
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param bot 机器人实例
// @param messageService 消息服务
func (b *botServiceImpl) RegisterHandlers(bot *telebot.Bot, messageService MessageService) {
    // 处理文本消息
    bot.Handle(telebot.OnText, func(c telebot.Context) error {
        text := c.Text()
        if len(text) == 0 {
            return nil
        }
        
        // 如果是命令，不处理
        if text[0] == '/' {
            return nil
        }
        
        // 如果文本较短，视为搜索查询
        if len([]rune(text)) <= 10 {
            return messageService.SearchWithPagination(c, text, 1, "message")
        }
        
        // 否则保存消息
        data := map[string]interface{}{
            "message_id": c.Message().ID,
            "chat_id":    c.Chat().ID,
            "chat_title": c.Chat().Title,
            "text":       text,
            "sender_id":  fmt.Sprintf("user_%d", c.Sender().ID),
            "date":       c.Message().Time().Format("2006-01-02T15:04:05Z07:00"),
        }
        
        return messageService.SaveMessage(data)
    })

    // 处理回调查询
    bot.Handle(telebot.OnCallback, messageService.HandleCallback)

    // 帮助命令
    bot.Handle("/help", func(c telebot.Context) error {
        helpText := `Available commands:
        /help - Show this help message
        /clong - Clone the bot
        /sponsor - Support the bot
        /mini - Enter mini mode
        Send a query (10 chars or less) to search with pagination and filters.`
        return c.Send(helpText)
    })

    // 克隆命令
    bot.Handle("/clong", func(c telebot.Context) error {
        cloneGuide := `To clone this bot:
1. Create a new bot via @BotFather on Telegram and get a Bot Token.
2. Clone the repository: https://github.com/your-repo/telegram-bot
3. Set environment variables (BOT_TOKEN, POCKETBASE_TOKEN, MEILISEARCH_KEY).
4. Deploy using Docker: 'docker run -e BOT_TOKEN=your_token your-image'
Visit https://your-repo.com for details.`
        return c.Send(cloneGuide)
    })

    // 赞助命令
    bot.Handle("/sponsor", func(c telebot.Context) error {
        return c.Send("Support us via: https://your-sponsor-link\nThank you!")
    })

    // 迷你模式命令
    bot.Handle("/mini", func(c telebot.Context) error {
        return c.Send("Mini mode activated! Send a query to search.")
    })
}

/*
 * 关键算法说明：
 * 1. 机器人管理：使用线程安全的map管理多个机器人实例
 * 2. Webhook处理：接收和处理Telegram Webhook更新
 * 3. 消息路由：根据消息类型和内容路由到不同的处理函数
 * 
 * 待优化事项：
 * 1. 动态配置：支持动态添加和移除机器人
 * 2. 状态监控：添加机器人状态监控和健康检查
 * 3. 错误处理：改进错误处理和恢复机制
 * 
 * 兼容性说明：
 * 1. 依赖telebot.v3库
 * 2. 需要有效的Telegram Bot Token
 * 3. 需要配置正确的Webhook URL
 */