package handler

import (
	"bot-service/internal/config"
	"bot-service/internal/index"
	"bot-service/internal/user"
	"bot-service/internal/usecase"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"gopkg.in/telebot.v4"
)

type GetChatMemberCountResponse struct {
	Ok     bool `json:"ok"`
	Result int  `json:"result"`
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

// BotHandler 定义机器人处理器接口
type BotHandler interface {
	InitBot(config BotConfig, fullConfig *config.Config) (*telebot.Bot, error)
	GetBot(token string) (*telebot.Bot, bool)
	ProcessUpdate(token string, update *telebot.Update) error
	RegisterHandlers(bot *telebot.Bot)
}

// botHandlerImpl 实现 BotHandler 接口
type botHandlerImpl struct {
	bots           map[string]*telebot.Bot
	mutex          sync.RWMutex
	messageUsecase usecase.MessageUsecase
	cfg            *config.Config
}

// NewBotHandler 创建新的机器人处理器实例
func NewBotHandler(messageUsecase usecase.MessageUsecase, cfg *config.Config) BotHandler {
	return &botHandlerImpl{
		bots:           make(map[string]*telebot.Bot),
		messageUsecase: messageUsecase,
		cfg:            cfg,
	}
}

// InitBot 初始化机器人
func (b *botHandlerImpl) InitBot(botConfig BotConfig, cfg *config.Config) (*telebot.Bot, error) {
	bot, err := telebot.NewBot(telebot.Settings{
		Token: botConfig.Token,
		URL:   cfg.Bot.APIEndpoint,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to init bot %s: %v", botConfig.Token, err)
	}

	b.mutex.Lock()
	b.bots[botConfig.Token] = bot
	b.mutex.Unlock()

	fmt.Printf("Bot %v initialized\n", bot)
	listenURL := cfg.Bot.WebhookURL + "/webhook?token=" + botConfig.Token
	fmt.Printf("Setting webhook %s\n", listenURL)

	webhook := &telebot.Webhook{Endpoint: &telebot.WebhookEndpoint{PublicURL: listenURL}}
	if err := bot.SetWebhook(webhook); err != nil {
		return nil, fmt.Errorf("failed to set webhook for bot %s: %v", botConfig.Token, err)
	}

	return bot, nil
}

// GetBot 获取机器人实例
func (b *botHandlerImpl) GetBot(token string) (*telebot.Bot, bool) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	bot, exists := b.bots[token]
	return bot, exists
}

// ProcessUpdate 处理Webhook更新
func (b *botHandlerImpl) ProcessUpdate(token string, update *telebot.Update) error {
	bot, exists := b.GetBot(token)
	if !exists {
		return fmt.Errorf("bot not found: %s", token)
	}
	bot.ProcessUpdate(*update)
	return nil
}

// getChatMemberCount 使用对 Telegram Bot API 的直接 HTTP 调用来检索聊天中的成员数。
func getChatMemberCount(token string, chatID int64) (int, error) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/getChatMemberCount?chat_id=%d", token, chatID)
	resp, err := http.Get(apiURL)
	if err != nil {
		slog.Error("获取聊天成员数量失败", "error", err)
		return 0, fmt.Errorf("failed to get chat member count: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("读取响应体失败", "error", err)
		return 0, fmt.Errorf("failed to read response body: %w", err)
	}

	var chatMemberCountResp GetChatMemberCountResponse
	if err := json.Unmarshal(body, &chatMemberCountResp); err != nil {
		slog.Error("解析 JSON 响应失败", "error", err)
		return 0, fmt.Errorf("failed to unmarshal json response: %w", err)
	}

	if !chatMemberCountResp.Ok {
		slog.Error("API 响应不成功", "body", string(body))
		return 0, fmt.Errorf("telegram API error: %s", string(body))
	}
	// 打印日志
	slog.Info("获取聊天成员数量成功", "chat_id", chatID, "member_count", chatMemberCountResp.Result)
	return chatMemberCountResp.Result, nil
}

// RegisterHandlers 注册消息处理函数
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param bot 机器人实例
func (b *botHandlerImpl) RegisterHandlers(bot *telebot.Bot) {
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
				chat, err := bot.ChatByUsername("@" + username)
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
				memberCount, err := getChatMemberCount(bot.Token, chat.ID)
				if err != nil {
					slog.Error("获取成员数量失败", "error", err, "chat", chat.Username)

					// 发送带有按钮的错误消息
					message := "获取用户数量失败，请将机器人拉入群组后重试。"
					inlineKeys := [][]telebot.InlineButton{
						{
							{
								Text: "🔄 重新获取",
								Data: "retry_index:" + text, // text 是原始的 https://t.me/... 链接
							},
							{
								Text: "➕ 添加到群组/频道",
								URL:  fmt.Sprintf("https://t.me/%s?startgroup=true", bot.Me.Username),
							},
						},
					}
					return c.Send(message, &telebot.SendOptions{
						ReplyMarkup: &telebot.ReplyMarkup{InlineKeyboard: inlineKeys},
					})
				}

				// 如果成功，继续索引
				data := map[string]interface{}{
					"chat_id":       fmt.Sprintf("%d", chat.ID),
					"type":          string(chat.Type),
					"title":         chat.Title,
					"username":      chat.Username,
					"first_name":    chat.FirstName,
					"last_name":     chat.LastName,
					"description":   description,
					"is_verified":   false,
					"members_count": memberCount,
					"created_at":    time.Now().Format("2006-01-02T15:04:05Z07:00"),
					"updated_at":    time.Now().Format("2006-01-02T15:04:05Z07:00"),
					"invite_link":   fullChat.InviteLink,
				}
				fmt.Printf("Data to save: %+v\n", data)
				if err := index.SaveTelegramIndex(b.cfg, data); err != nil {
					fmt.Printf("Error saving index: %v\n", err)
					return err
				}
				fmt.Println("Index saved successfully")

				// 构建并发送成功消息
				successMessage := fmt.Sprintf(
					"<b>群组收录成功</b>\n\n"+
						"<b>标题:</b> %s\n"+
						"<b>用户名:</b> @%s\n"+
						"<b>描述:</b> %s\n"+
						"<b>成员数量:</b> %d",
					html.EscapeString(chat.Title),
					html.EscapeString(chat.Username),
					html.EscapeString(description),
					memberCount,
				)
				return c.Send(successMessage, &telebot.SendOptions{ParseMode: telebot.ModeHTML})
			}
		}

		// 如果文本长度小于10，则触发搜索
		if utf8.RuneCountInString(text) < 10 {
			return b.messageUsecase.SearchWithPagination(c, text, 1, "")
		}

		// 默认保存消息
		// 首先，为 operation_details 字段创建一个详细信息映射
		details := map[string]interface{}{
			"message_id":    c.Message().ID,
			"chat_id":       c.Chat().ID,
			"chat_title":    c.Chat().Title,
			"chat_username": c.Chat().Username,
			"chat_type":     string(c.Chat().Type),
			"text":          text,
			"sender_is_bot": c.Sender().IsBot,
			"date":          c.Message().Time().Format("2006-01-02T15:04:05Z07:00"),
		}

		detailsJSON, err := json.Marshal(details)
		if err != nil {
			// 记录错误，但不阻止用户
			fmt.Printf("ERROR: Failed to marshal operation details: %v\n", err)
		}

		data := map[string]interface{}{
			"user":              c.Sender().ID,
			"bot_id":            fmt.Sprintf("%d", bot.Me.ID),
			"operation_type":    "save_message",
			"operation_time":    time.Now().Format("2006-01-02T15:04:05Z07:00"),
			"operation_details": string(detailsJSON),
			"create_time":       time.Now().Format("2006-01-02T15:04:05Z07:00"),
		}

		return b.messageUsecase.SaveMessage(data)
	})

	// 处理回调查询
	bot.Handle(telebot.OnCallback, b.messageUsecase.HandleCallback)

	// /search 命令处理
	bot.Handle("/search", func(c telebot.Context) error {
		query := c.Message().Payload
		if query == "" {
			return c.Send("Please provide a search query. Usage: /search <query>")
		}
		return b.messageUsecase.SearchWithPagination(c, query, 1, "")
	})

	// /start 命令处理
	bot.Handle("/start", func(c telebot.Context) error {
		// 保存或更新用户信息
		if err := user.SaveUser(c.Sender()); err != nil {
			// 记录错误，但不影响用户体验
			fmt.Printf("保存用户信息失败: %v\n", err)
		}

		// 欢迎消息
		welcomeText := fmt.Sprintf(
			"你好, %s, 欢迎来到我们的TG机器人！\n\n"+
				"<a href=\"https://t.me/addlist/pMIbwEotf14wOGU1\">👏 点击加入我们的交流大群 👏</a>\n\n"+
				"<b>使用说明:</b>\n"+
				"- 直接向我发送消息，即可将内容保存到您的个人收藏夹。\n"+
				"- 发送短于10个字符的文本，将触发搜索功能。\n"+
				"- 使用 <code>/mini</code> 命令可以随时唤出小程序。\n\n"+
				"🔍✨👇 点击下方按钮打开小程序，或选择一个大群加入我们！",
			c.Sender().FirstName,
		)

		// 创建内联键盘
		inlineKeys := [][]telebot.InlineButton{
			{},
			{
				telebot.InlineButton{Text: "搜索大群", URL: "https://t.me/SoSo00000000001"},
				telebot.InlineButton{Text: "搜索每日更新频道", URL: "https://t.me/SoSo00000000002"},
			},
			{
				telebot.InlineButton{Text: "搜索消息监听", URL: "https://t.me/SoSo00000000003"},
				telebot.InlineButton{
					Text: "🚀 打开小程序",
					WebApp: &telebot.WebApp{
						URL: "https://timi000000001-ai.github.io/index-telegram-app", // TODO: URL应可配置
					},
				},
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
		helpText := `<b>可用命令列表：</b>

/help - 显示此帮助信息
/search <关键词> - 搜索群组、频道和消息
/clong - 克隆机器人
/sponsor - 支持我们
/mini - 打开小程序
/disclaimer - 查看免责声明

<b>使用说明：</b>
1. 直接发送消息给机器人，消息会被保存到您的个人收藏夹
2. 使用 /search 命令搜索群组、频道和消息
3. 搜索结果支持分页和过滤功能
4. 点击搜索结果中的链接可以直接访问`

		return c.Send(helpText, &telebot.SendOptions{
			ParseMode: telebot.ModeHTML,
		})
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
