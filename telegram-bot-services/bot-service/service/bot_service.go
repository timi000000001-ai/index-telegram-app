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
	"bot-service/internal/config"
	"bot-service/internal/index"
	"bot-service/internal/user"
	"fmt"
	"strings"
	"sync"
	"time"

	"gopkg.in/telebot.v3"
)

var (
	defaultBotService *botServiceImpl
	once              sync.Once
)

// Init initializes the default bot service.
func Init(botConfig BotConfig) (*telebot.Bot, error) {
	var err error
	var bot *telebot.Bot
	once.Do(func() {
		cfg, loadErr := config.LoadConfig("development")
		if loadErr != nil {
			err = fmt.Errorf("failed to load config: %w", loadErr)
			return
		}

		storageService := NewStorageService(StorageConfig{
			PocketBaseURL:    cfg.Storage.PocketBaseURL,
			MeilisearchURL:   cfg.Storage.MeilisearchURL,
			MeilisearchToken: cfg.Storage.MeilisearchToken,
		})
		searchService := NewSearchService(SearchConfig{
			MeilisearchURL:       cfg.Search.MeilisearchURL,
			MeilisearchKey:       cfg.Search.MeilisearchKey,
			ManagementServiceURL: cfg.Search.ManagementServiceURL,
		})
		messageService := NewMessageService(storageService, searchService)

		defaultBotService = &botServiceImpl{
			bots: make(map[string]*telebot.Bot),
		}

		var initErr error
		bot, initErr = defaultBotService.initBot(botConfig, cfg)
		if initErr != nil {
			err = initErr
			return
		}
		defaultBotService.registerHandlers(bot, messageService)
	})
	if err != nil {
		return nil, err
	}
	if bot == nil {
		bot, _ = defaultBotService.getBot(botConfig.Token)
	}
	return bot, nil
}

// ProcessUpdate processes a webhook update.
func ProcessUpdate(token string, update *telebot.Update) error {
	if defaultBotService == nil {
		return fmt.Errorf("service not initialized")
	}
	return defaultBotService.processUpdate(token, update)
}

// BotConfig 定义机器人配置
// @author fcj
// @date 2023-11-15
// @version 1.0.0
type BotConfig struct {
	Token                  string `json:"bot_token"`
	Name                   string `json:"bot_name"`
	Status                 string `json:"bot_status"`
	WebhookURL             string `json:"webhook_url"`
	ManagementServiceURL   string `json:"management_service_url"`
	ManagementServiceToken string `json:"management_service_token"`
}

// BotService 定义机器人服务接口
// @author fcj
// @date 2023-11-15
// @version 1.0.0
type BotService interface {
	// initBot 初始化机器人
	initBot(config BotConfig, fullConfig *config.Config) (*telebot.Bot, error)

	// getBot 获取机器人实例
	getBot(token string) (*telebot.Bot, bool)

	// processUpdate 处理Webhook更新
	processUpdate(token string, update *telebot.Update) error

	// registerHandlers 注册消息处理函数
	registerHandlers(bot *telebot.Bot, messageService MessageService)
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

// initBot 初始化机器人
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param config 机器人配置
// @return error 错误信息
func (b *botServiceImpl) initBot(botConfig BotConfig, cfg *config.Config) (*telebot.Bot, error) {
	bot, err := telebot.NewBot(telebot.Settings{
		Token: botConfig.Token,
		URL:   cfg.Bot.APIEndpoint, // 仅用于本地机器人API服务器
	})

	if err != nil {
		return nil, fmt.Errorf("failed to init bot %s: %v", botConfig.Token, err)
	}

	b.mutex.Lock()
	b.bots[botConfig.Token] = bot
	b.mutex.Unlock()

	//打印机器人信息
	fmt.Printf("Bot %v initialized\n", bot)
	// Set webhook. This only registers the URL with Telegram.
	// The actual HTTP server that listens for updates is started in main.go.
	listenUrl := cfg.Bot.WebhookURL + "/webhook?token=" + botConfig.Token
	fmt.Printf("Setting webhook %s\n", listenUrl)

	webhook := &telebot.Webhook{Endpoint: &telebot.WebhookEndpoint{PublicURL: listenUrl}}
	//打印日志
	if err := bot.SetWebhook(webhook); err != nil {
		return nil, fmt.Errorf("failed to set webhook for bot %s: %v", botConfig.Token, err)
	}

	return bot, nil
}

// getBot 获取机器人实例
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param token 机器人令牌
// @return *telebot.Bot 机器人实例
// @return bool 是否存在
func (b *botServiceImpl) getBot(token string) (*telebot.Bot, bool) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	bot, exists := b.bots[token]
	return bot, exists
}

// processUpdate 处理Webhook更新
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param token 机器人令牌
// @param update Webhook更新
// @return error 错误信息
func (b *botServiceImpl) processUpdate(token string, update *telebot.Update) error {
	bot, exists := b.getBot(token)
	if !exists {
		return fmt.Errorf("bot not found: %s", token)
	}

	bot.ProcessUpdate(*update)
	return nil
}

// registerHandlers 注册消息处理函数
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param bot 机器人实例
// @param messageService 消息服务
func (b *botServiceImpl) registerHandlers(bot *telebot.Bot, messageService MessageService) {
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

		// 如果文本以 "https://t.me" 开头，处理为索引请求
		if strings.HasPrefix(text, "https://t.me/") {
				fmt.Printf("Processing link: %s\n", text)
				// 解析用户名（假设链接为 https://t.me/username）
				parts := strings.Split(strings.TrimPrefix(text, "https://t.me/"), "/")
				if len(parts) > 0 {
					username := parts[0]
					fmt.Printf("Username: %s\n", username)
					chat, err := bot.ChatByUsername("@"+username)
					if err != nil {
						fmt.Printf("Error getting chat: %v\n", err)
						return err
					}
					fmt.Printf("Chat retrieved: %+v\n", chat)
					// 获取完整聊天信息
					fullChat, err := bot.ChatByID(chat.ID)
					if err != nil {
						fmt.Printf("Error getting full chat: %v\n", err)
						fullChat = chat
					}
					fmt.Printf("Full chat: %+v\n", fullChat)
					description := fullChat.Description
					if chat.Type == telebot.ChatPrivate {
						description = fullChat.Bio
					}
					data := map[string]interface{}{
						"chat_id": fmt.Sprintf("%d", chat.ID),
						"type": string(chat.Type),
						"title": chat.Title,
						"username": chat.Username,
						"first_name": chat.FirstName,
						"last_name": chat.LastName,
						"description": description,
						"is_verified": false, // 默认值，因为 Verified 未定义
						"members_count": 0, // 默认值，因为 MemberCount 未定义；可添加单独查询
						"created_at": time.Now().Format("2006-01-02T15:04:05Z07:00"),
						"updated_at": time.Now().Format("2006-01-02T15:04:05Z07:00"),
						"invite_link": fullChat.InviteLink,
						// 添加更多可用字段
					}
					fmt.Printf("Data to save: %+v\n", data)
					if err := index.SaveTelegramIndex(data); err != nil {
						fmt.Printf("Error saving index: %v\n", err)
						return err
					}
					fmt.Println("Index saved successfully")
					return c.Send("已索引聊天信息")
				}
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

	// /start 命令处理
	bot.Handle("/start", func(c telebot.Context) error {
		// 保存或更新用户信息
		if err := user.SaveUser(c.Sender()); err != nil {
			// 记录错误，但不影响用户体验
			fmt.Printf("保存用户信息失败: %v\n", err)
		}

		// 欢迎消息
		welcomeText := fmt.Sprintf(
			"你好, %s, 欢迎来到我们的TG机器人！\n\n" +
				"<a href=\"https://t.me/addlist/pMIbwEotf14wOGU1\">👏 点击加入我们的交流大群 👏</a>\n\n" +
				"<b>使用说明:</b>\n" +
				"- 直接向我发送消息，即可将内容保存到您的个人收藏夹。\n" +
				"- 发送短于10个字符的文本，将触发搜索功能。\n" +
				"- 使用 <code>/mini</code> 命令可以随时唤出小程序。\n\n" +
				"🔍✨👇 点击下方按钮打开小程序，或选择一个大群加入我们！",
			c.Sender().FirstName,
		)

		// 创建内联键盘
		inlineKeys := [][]telebot.InlineButton{
			{
				telebot.InlineButton{
					Text: "🚀 打开小程序",
					WebApp: &telebot.WebApp{
						URL: "https://timi000000001-ai.github.io/index-telegram-app", // TODO: URL应可配置
					},
				},
			},
			{
				telebot.InlineButton{Text: "中文学习交流群", URL: "https://t.me/addlist/pMIbwEotf14wOGU1"},
				telebot.InlineButton{Text: "资源分享群", URL: "https://t.me/addlist/pMIbwEotf14wOGU1"},
			},
			{
				telebot.InlineButton{Text: "技术交流群", URL: "https://t.me/addlist/pMIbwEotf14wOGU1"},
				telebot.InlineButton{Text: "闲聊吹水群", URL: "https://t.me/addlist/pMIbwEotf14wOGU1"},
			},
		}

		return c.Send(welcomeText, &telebot.SendOptions{
			ParseMode:             telebot.ModeHTML,
			DisableWebPagePreview: true,
			ReplyMarkup: &telebot.ReplyMarkup{
				InlineKeyboard: inlineKeys,
			},
		})
	})

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
		sponsorText := "如果您觉得本机器人对您有帮助，请考虑赞助我们。\n\nTRX & USDT (TRC20):\n\n✨<code>TD5JGaR7cY5ZxDnZNgmCSv66axR9DhrcYz</code>✨\n\n"

		inlineKeys := [][]telebot.InlineButton{
			{
				telebot.InlineButton{Text: "联系我们", URL: "https://t.me/simi001001"},
				telebot.InlineButton{
					Text: "🚀 打开小程序",
					WebApp: &telebot.WebApp{
						URL: "https://timi000000001-ai.github.io/index-telegram-app", // TODO: URL应可配置
					},
				},
			},
		}

		return c.Send(sponsorText, &telebot.SendOptions{
			ParseMode: telebot.ModeHTML,
			ReplyMarkup: &telebot.ReplyMarkup{
				InlineKeyboard: inlineKeys,
			},
		})
	})

	// 迷你模式命令
	bot.Handle("/mini", func(c telebot.Context) error {
		// 创建一个内联键盘，其中包含一个小程序按钮
		inlineKeys := [][]telebot.InlineButton{
			{
				telebot.InlineButton{
					Text: "打开小程序",
					WebApp: &telebot.WebApp{
						URL: "https://timi000000001-ai.github.io/index-telegram-app", // TODO: Make this URL configurable
					},
				},
			},
		}

		// 发送带有内联键盘的消息
		return c.Send("<a href=\"https://t.me/addlist/pMIbwEotf14wOGU1\">👏 加入搜索大群 👏 \n\n💡 温馨提示：加入群组可获取更多优质资源！</a>  \n\n🔍✨👇 点击下方按钮打开小程序 🚀\n", &telebot.SendOptions{
			ParseMode:             telebot.ModeHTML,
			DisableWebPagePreview: true,
			ReplyMarkup: &telebot.ReplyMarkup{
				InlineKeyboard: inlineKeys,
			},
		})
	})

	// 免责声明命令
	bot.Handle("/disclaimer", func(c telebot.Context) error {
		disclaimerText := `⚠️ <b>法律声明</b> ⚠️

<b>使用限制</b>：本项目不适用于中国大陆。Telegram 在中国大陆受到政府的访问限制，本项目的数据收集和处理活动可能违反当地法律法规。

<b>免责声明</b>：本项目开发人员对因使用不当、违反当地法律或数据隐私问题而导致的任何后果概不负责。用户应自行评估法律风险，并在必要时咨询法律专业人士。

<b>建议</b>：如果您位于中国大陆，请不要下载、安装或运行本项目。请寻找符合当地法规的替代方案。`
		return c.Send(disclaimerText, &telebot.SendOptions{
			ParseMode: telebot.ModeHTML,
		})
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
