/*
 * 文件功能描述：集合服务，处理消息采集的核心逻辑
 * 主要类/接口说明：CollectionService接口及其实现
 * 修改历史记录：
 * @author fcj
 * @date 2023-11-15
 * @version 1.0.0
 * © Telegram Bot Services Team
 */

package service

import (
    "fmt"
    "log"
    "time"
)

// CollectionConfig 采集配置
// @author fcj
// @date 2023-11-15
// @version 1.0.0
type SearchConfig struct {
    MeilisearchURL  string `json:"meilisearch_url"`
    MeilisearchToken string `json:"meilisearch_token"`
    MessageLimit     int    `json:"message_limit"`
}

// CollectionService 定义消息采集服务接口
// @author fcj
// @date 2023-11-15
// @version 1.0.0
type CollectionService interface {
    // StartCollection 开始采集消息
    StartCollection() error
    
    // CollectFromChat 从特定群组采集消息
    CollectFromChat(phoneNumber string, chatID int64) error
    
    // CollectFromAllChats 从所有配置的群组采集消息
    CollectFromAllChats() error
    
    // ScheduleCollection 定时采集消息
    ScheduleCollection(interval time.Duration) error
    
    // StopCollection 停止采集消息
    StopCollection() error
}

// collectionServiceImpl 实现CollectionService接口
// @author fcj
// @date 2023-11-15
// @version 1.0.0
type collectionServiceImpl struct {
    config          SearchConfig
    telegramService TelegramService
    sessionService  SessionService
    messageLimit    int
    ticker          *time.Ticker
    stopChan        chan struct{}
    isRunning       bool
}

// NewCollectionService 创建新的消息采集服务实例
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param config 采集配置
// @param telegramService Telegram客户端服务
// @param sessionService 会话管理服务
// @return CollectionService 消息采集服务实例
func NewCollectionService(config SearchConfig, telegramService TelegramService, sessionService SessionService) CollectionService {
    return &collectionServiceImpl{
        config:          config,
        telegramService: telegramService,
        sessionService:  sessionService,
        messageLimit:    config.MessageLimit,
        stopChan:        make(chan struct{}),
        isRunning:       false,
    }
}

// StartCollection 开始采集消息
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @return error 错误信息
func (c *collectionServiceImpl) StartCollection() error {
    // 获取所有配置的群组ID
    chatIDs := c.sessionService.GetConfiguredChatIDs()
    if len(chatIDs) == 0 {
        return fmt.Errorf("no chat IDs configured for collection")
    }
    
    // 获取所有活跃会话
    activeSessions := c.sessionService.ListActiveSessions()
    if len(activeSessions) == 0 {
        return fmt.Errorf("no active sessions available for collection")
    }
    
    log.Printf("Starting collection for %d chats using %d active sessions", len(chatIDs), len(activeSessions))
    
    // 可以优化为负载均衡或轮询多个会话
    // 目前直接使用 telegramService 进行采集
    
    // 采集所有配置的群组消息
    return c.telegramService.CollectAllConfiguredMessages(chatIDs)
}

// CollectFromChat 从特定群组采集消息
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param phoneNumber 电话号码
// @param chatID 群组ID
// @return error 错误信息
func (c *collectionServiceImpl) CollectFromChat(phoneNumber string, chatID int64) error {
    // 检查会话是否存在
    _, err := c.sessionService.GetSession(phoneNumber)
    if err != nil {
        return fmt.Errorf("session not found: %w", err)
    }
    
    // 采集消息
    return c.telegramService.CollectMessages(phoneNumber, chatID, c.messageLimit)
}

// CollectFromAllChats 从所有配置的群组采集消息
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @return error 错误信息
func (c *collectionServiceImpl) CollectFromAllChats() error {
    // 获取所有配置的群组ID
    chatIDs := c.sessionService.GetConfiguredChatIDs()
    if len(chatIDs) == 0 {
        return fmt.Errorf("no chat IDs configured for collection")
    }
    
    // 采集所有配置的群组消息
    return c.telegramService.CollectAllConfiguredMessages(chatIDs)
}

// ScheduleCollection 定时采集消息
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param interval 采集间隔
// @return error 错误信息
func (c *collectionServiceImpl) ScheduleCollection(interval time.Duration) error {
    if c.isRunning {
        return fmt.Errorf("collection is already scheduled")
    }
    
    // 创建定时器
    c.ticker = time.NewTicker(interval)
    c.isRunning = true
    
    // 启动定时采集
    go func() {
        for {
            select {
            case <-c.ticker.C:
                // 检查是否有配置的群组和活跃会话
                chatIDs := c.sessionService.GetConfiguredChatIDs()
                activeSessions := c.sessionService.ListActiveSessions()
                
                if len(chatIDs) > 0 && len(activeSessions) > 0 {
                    log.Printf("Starting scheduled collection for %d chats", len(chatIDs))
                    if err := c.StartCollection(); err != nil {
                        log.Printf("Scheduled collection failed: %v", err)
                    }
                } else {
                    log.Printf("Skipping scheduled collection: no chats (%d) or active sessions (%d)", len(chatIDs), len(activeSessions))
                }
            case <-c.stopChan:
                c.ticker.Stop()
                c.isRunning = false
                return
            }
        }
    }()
    
    log.Printf("Scheduled collection every %v", interval)
    return nil
}

// StopCollection 停止采集消息
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @return error 错误信息
func (c *collectionServiceImpl) StopCollection() error {
    if !c.isRunning {
        return fmt.Errorf("collection is not running")
    }
    
    // 停止定时采集
    c.stopChan <- struct{}{}
    log.Printf("Stopped scheduled collection")
    return nil
}

/*
 * 关键算法说明：
 * 1. 定时采集：使用time.Ticker实现定时采集
 * 2. 会话管理：通过SessionService获取活跃会话
 * 3. 消息采集：通过TelegramService采集群组消息
 * 
 * 待优化事项：
 * 1. 负载均衡：实现多会话负载均衡采集
 * 2. 错误重试：添加采集失败重试机制
 * 3. 增量采集：实现增量采集以减少重复数据
 * 
 * 兼容性说明：
 * 1. 依赖TelegramService和SessionService
 * 2. 支持自定义采集间隔和消息数量限制
 */