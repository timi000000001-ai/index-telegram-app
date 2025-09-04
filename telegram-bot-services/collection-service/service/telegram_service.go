/*
 * 文件功能描述：Telegram客户端服务，处理Telegram API交互
 * 主要类/接口说明：TelegramService接口及其实现
 * 修改历史记录：
 * @author fcj
 * @date 2023-11-15
 * @version 1.0.0
 * © Telegram Bot Services Team
 */

package service

import (
    "context"
    "fmt"
    "log"
    "sync"
    "time"

    "github.com/gotd/td/telegram/auth"
    tdapi "github.com/gotd/td/telegram"
    "github.com/gotd/td/tg"
)

// TelegramConfig Telegram客户端配置
// @author fcj
// @date 2023-11-15
// @version 1.0.0
type TelegramConfig struct {
    AppID   int    `json:"app_id"`
    AppHash string `json:"app_hash"`
}

// SessionData 会话数据结构
// @author fcj
// @date 2023-11-15
// @version 1.0.0
type SessionData struct {
    Client *tdapi.Client
}

// TelegramService 定义Telegram客户端服务接口
// @author fcj
// @date 2023-11-15
// @version 1.0.0
type TelegramService interface {
    // NewClient 创建新的Telegram客户端
    NewClient(phoneNumber string) (*tdapi.Client, error)
    
    // SendCode 发送验证码
    SendCode(phoneNumber string) error
    
    // Login 登录Telegram账号
    Login(phoneNumber string, code string) error
    
    // IsAuthorized 检查客户端是否已认证
    IsAuthorized(phoneNumber string) (bool, error)
    
    // CollectMessages 采集群组消息
    CollectMessages(phoneNumber string, chatID int64, limit int) error
    
    // CollectAllConfiguredMessages 采集所有配置的群组消息
    CollectAllConfiguredMessages(chatIDs []int64) error
}

// telegramServiceImpl 实现TelegramService接口
// @author fcj
// @date 2023-11-15
// @version 1.0.0
type telegramServiceImpl struct {
    config      TelegramConfig
    sessions    sync.Map
    storageService StorageService
}

// NewTelegramService 创建新的Telegram客户端服务实例
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param config Telegram客户端配置
// @param storageService 存储服务
// @return TelegramService Telegram客户端服务实例
func NewTelegramService(config TelegramConfig, storageService StorageService) TelegramService {
    return &telegramServiceImpl{
        config:      config,
        sessions:    sync.Map{},
        storageService: storageService,
    }
}

// NewClient 创建新的Telegram客户端
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param phoneNumber 电话号码
// @return *tdapi.Client Telegram客户端实例
// @return error 错误信息
func (t *telegramServiceImpl) NewClient(phoneNumber string) (*tdapi.Client, error) {
    if phoneNumber == "" {
        return nil, fmt.Errorf("phone number cannot be empty")
    }
    
    // 创建新的Telegram客户端
    client := tdapi.NewClient(t.config.AppID, t.config.AppHash, tdapi.Options{})
    return client, nil
}

// SendCode 发送验证码
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param phoneNumber 电话号码
// @return error 错误信息
func (t *telegramServiceImpl) SendCode(phoneNumber string) error {
    // 检查是否已经有会话
    sessionInterface, exists := t.sessions.Load(phoneNumber)
    
    // 创建或获取客户端
    var client *tdapi.Client
    
    if exists {
        // 获取客户端实例
        switch v := sessionInterface.(type) {
        case *tdapi.Client:
            client = v
        case SessionData:
            client = v.Client
        default:
            return fmt.Errorf("invalid session data type for phone %s: %T", phoneNumber, sessionInterface)
        }
        
        // 检查客户端是否已认证
        isAuthorized, err := t.IsAuthorized(phoneNumber)
        if err == nil && isAuthorized {
            return fmt.Errorf("already logged in")
        }
    } else {
        // 创建新客户端
        var err error
        client, err = t.NewClient(phoneNumber)
        if err != nil {
            return fmt.Errorf("failed to create client: %w", err)
        }
    }
    
    // 发送验证码
    err := client.Run(context.Background(), func(ctx context.Context) error {
        _, err := client.Auth().SendCode(ctx, phoneNumber, auth.SendCodeOptions{})
        if err != nil {
            return fmt.Errorf("failed to send verification code: %w", err)
        }
        return nil
    })
    
    if err != nil {
        return err
    }
    
    // 存储客户端
    sessionData := SessionData{
        Client: client,
    }
    t.sessions.Store(phoneNumber, sessionData)
    return nil
}

// Login 登录Telegram账号
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param phoneNumber 电话号码
// @param code 验证码
// @return error 错误信息
func (t *telegramServiceImpl) Login(phoneNumber string, code string) error {
    // 获取会话
    sessionInterface, exists := t.sessions.Load(phoneNumber)
    if !exists {
        return fmt.Errorf("no session found for phone %s", phoneNumber)
    }
    
    // 获取客户端实例
    var client *tdapi.Client
    switch v := sessionInterface.(type) {
    case *tdapi.Client:
        client = v
    case SessionData:
        client = v.Client
    default:
        return fmt.Errorf("invalid session data type for phone %s: %T", phoneNumber, sessionInterface)
    }
    
    if client == nil {
        return fmt.Errorf("no valid client found for phone %s", phoneNumber)
    }
    
    // 创建认证流程
    codeAuth := auth.CodeAuthenticatorFunc(func(ctx context.Context, sentCode *tg.AuthSentCode) (string, error) {
        return code, nil
    })
    
    flow := auth.NewFlow(
        auth.CodeOnly(phoneNumber, codeAuth),
        auth.SendCodeOptions{},
    )
    
    // 运行认证流程
    err := client.Run(context.Background(), func(ctx context.Context) error {
        return flow.Run(ctx, client.Auth())
    })
    
    if err != nil {
        return fmt.Errorf("login failed: %w", err)
    }
    
    // 更新会话
    t.sessions.Store(phoneNumber, client)
    return nil
}

// IsAuthorized 检查客户端是否已认证
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param phoneNumber 电话号码
// @return bool 是否已认证
// @return error 错误信息
func (t *telegramServiceImpl) IsAuthorized(phoneNumber string) (bool, error) {
    // 获取会话
    sessionInterface, exists := t.sessions.Load(phoneNumber)
    if !exists {
        return false, fmt.Errorf("no session found for phone %s", phoneNumber)
    }
    
    // 获取客户端实例
    var client *tdapi.Client
    switch v := sessionInterface.(type) {
    case *tdapi.Client:
        client = v
    case SessionData:
        client = v.Client
    default:
        return false, fmt.Errorf("invalid session data type for phone %s: %T", phoneNumber, sessionInterface)
    }
    
    if client == nil {
        return false, fmt.Errorf("no valid client found for phone %s", phoneNumber)
    }
    
    // 检查认证状态
    var isAuthorized bool
    err := client.Run(context.Background(), func(ctx context.Context) error {
        status, err := client.Auth().Status(ctx)
        if err != nil {
            return err
        }
        isAuthorized = status.Authorized
        return nil
    })
    
    return isAuthorized, err
}

// CollectMessages 采集群组消息
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param phoneNumber 电话号码
// @param chatID 群组ID
// @param limit 消息数量限制
// @return error 错误信息
func (t *telegramServiceImpl) CollectMessages(phoneNumber string, chatID int64, limit int) error {
    // 获取会话
    sessionInterface, exists := t.sessions.Load(phoneNumber)
    if !exists {
        return fmt.Errorf("no session found for phone %s", phoneNumber)
    }
    
    // 获取客户端实例
    var client *tdapi.Client
    switch v := sessionInterface.(type) {
    case *tdapi.Client:
        client = v
    case SessionData:
        client = v.Client
    default:
        return fmt.Errorf("invalid session data type for phone %s: %T", phoneNumber, sessionInterface)
    }
    
    if client == nil {
        return fmt.Errorf("no valid client found for phone %s", phoneNumber)
    }
    
    // 检查认证状态
    isAuthorized, err := t.IsAuthorized(phoneNumber)
    if err != nil || !isAuthorized {
        return fmt.Errorf("client for phone %s is not authorized: %v", phoneNumber, err)
    }
    
    // 采集消息
    return client.Run(context.Background(), func(ctx context.Context) error {
        api := client.API()
        
        // 构建请求获取历史消息
        messages, err := api.MessagesGetHistory(ctx, &tg.MessagesGetHistoryRequest{
            Peer: &tg.InputPeerChannel{
                ChannelID:  chatID,
                AccessHash: 0, // 需从管理端服务动态获取
            },
            Limit: limit,
        })
        if err != nil {
            return fmt.Errorf("failed to get message history: %w", err)
        }
        
        // 类型断言检查
        channelMessages, ok := messages.(*tg.MessagesChannelMessages)
        if !ok {
            return fmt.Errorf("unexpected message type: %T", messages)
        }
        
        log.Printf("Retrieved %d messages from chat %d", len(channelMessages.Messages), chatID)
        
        // 处理每条消息
        for _, msg := range channelMessages.Messages {
            message, ok := msg.(*tg.Message)
            if !ok {
                log.Printf("Skipping non-message type: %T", msg)
                continue
            }
            
            // 构建消息数据
            data := map[string]interface{}{
                "message_id": message.ID,
                "chat_id":    chatID,
                "chat_title": "Unknown", // 可通过 channels.GetFullChannel 获取
                "text":       message.Message,
                "sender_id":  fmt.Sprintf("user_%d", message.FromID),
                "date":       time.Unix(int64(message.Date), 0).Format(time.RFC3339),
            }
            
            // 保存到数据库和搜索引擎
            if err := t.storageService.SaveAndIndex(data); err != nil {
                log.Printf("Failed to save and index message %d: %v", message.ID, err)
            }
        }
        return nil
    })
}

// CollectAllConfiguredMessages 采集所有配置的群组消息
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param chatIDs 群组ID列表
// @return error 错误信息
func (t *telegramServiceImpl) CollectAllConfiguredMessages(chatIDs []int64) error {
    if len(chatIDs) == 0 {
        return fmt.Errorf("no chat IDs configured for collection")
    }
    
    // 遍历所有会话进行消息采集
    var lastError error
    t.sessions.Range(func(key, value interface{}) bool {
        phoneNumber, ok := key.(string)
        if !ok {
            log.Printf("Invalid session key type: %T", key)
            return true
        }
        
        // 检查认证状态
        isAuthorized, err := t.IsAuthorized(phoneNumber)
        if err != nil || !isAuthorized {
            log.Printf("Client for phone %s is not authorized: %v", phoneNumber, err)
            return true
        }
        
        log.Printf("Starting collection for phone %s", phoneNumber)
        for _, chatID := range chatIDs {
            // 默认采集 100 条消息
            if err := t.CollectMessages(phoneNumber, chatID, 100); err != nil {
                log.Printf("Failed to collect from chat %d: %v", chatID, err)
                lastError = err
            } else {
                log.Printf("Successfully collected messages from chat %d", chatID)
            }
        }
        return true
    })
    
    return lastError
}

/*
 * 关键算法说明：
 * 1. 会话管理：使用sync.Map存储和管理多个用户会话
 * 2. 认证流程：实现Telegram客户端的认证流程
 * 3. 消息采集：通过Telegram API获取群组历史消息
 * 
 * 待优化事项：
 * 1. 会话持久化：将会话数据持久化到数据库
 * 2. 错误重试：添加重试机制处理临时网络故障
 * 3. 并发控制：优化并发采集性能
 * 
 * 兼容性说明：
 * 1. 依赖gotd/td库与Telegram API交互
 * 2. 需要有效的Telegram API ID和Hash
 */