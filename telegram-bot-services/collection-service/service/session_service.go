/*
 * 文件功能描述：会话管理服务，处理用户会话状态
 * 主要类/接口说明：SessionService接口及其实现
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
    "time"
)

// Session 会话数据结构
// @author fcj
// @date 2023-11-15
// @version 1.0.0
type Session struct {
    PhoneNumber string
    ChatIDs     []int64
    LastActive  time.Time
}

// SessionService 定义会话管理服务接口
// @author fcj
// @date 2023-11-15
// @version 1.0.0
type SessionService interface {
    // CreateSession 创建新会话
    CreateSession(phoneNumber string) error
    
    // GetSession 获取会话
    GetSession(phoneNumber string) (*Session, error)
    
    // UpdateSession 更新会话
    UpdateSession(phoneNumber string, chatIDs []int64) error
    
    // DeleteSession 删除会话
    DeleteSession(phoneNumber string) error
    
    // ListActiveSessions 列出所有活跃会话
    ListActiveSessions() []*Session
    
    // GetConfiguredChatIDs 获取所有配置的群组ID
    GetConfiguredChatIDs() []int64
}

// sessionServiceImpl 实现SessionService接口
// @author fcj
// @date 2023-11-15
// @version 1.0.0
type sessionServiceImpl struct {
    sessions sync.Map
}

// NewSessionService 创建新的会话管理服务实例
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @return SessionService 会话管理服务实例
func NewSessionService() SessionService {
    return &sessionServiceImpl{
        sessions: sync.Map{},
    }
}

// CreateSession 创建新会话
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param phoneNumber 电话号码
// @return error 错误信息
func (s *sessionServiceImpl) CreateSession(phoneNumber string) error {
    if phoneNumber == "" {
        return fmt.Errorf("phone number cannot be empty")
    }
    
    // 检查是否已存在
    if _, exists := s.sessions.Load(phoneNumber); exists {
        return fmt.Errorf("session already exists for phone %s", phoneNumber)
    }
    
    // 创建新会话
    session := &Session{
        PhoneNumber: phoneNumber,
        ChatIDs:     []int64{},
        LastActive:  time.Now(),
    }
    
    s.sessions.Store(phoneNumber, session)
    return nil
}

// GetSession 获取会话
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param phoneNumber 电话号码
// @return *Session 会话
// @return error 错误信息
func (s *sessionServiceImpl) GetSession(phoneNumber string) (*Session, error) {
    if phoneNumber == "" {
        return nil, fmt.Errorf("phone number cannot be empty")
    }
    
    // 获取会话
    sessionInterface, exists := s.sessions.Load(phoneNumber)
    if !exists {
        return nil, fmt.Errorf("no session found for phone %s", phoneNumber)
    }
    
    // 类型断言
    session, ok := sessionInterface.(*Session)
    if !ok {
        return nil, fmt.Errorf("invalid session data type for phone %s: %T", phoneNumber, sessionInterface)
    }
    
    // 更新最后活跃时间
    session.LastActive = time.Now()
    s.sessions.Store(phoneNumber, session)
    
    return session, nil
}

// UpdateSession 更新会话
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param phoneNumber 电话号码
// @param chatIDs 群组ID列表
// @return error 错误信息
func (s *sessionServiceImpl) UpdateSession(phoneNumber string, chatIDs []int64) error {
    if phoneNumber == "" {
        return fmt.Errorf("phone number cannot be empty")
    }
    
    // 获取会话
    sessionInterface, exists := s.sessions.Load(phoneNumber)
    if !exists {
        return fmt.Errorf("no session found for phone %s", phoneNumber)
    }
    
    // 类型断言
    session, ok := sessionInterface.(*Session)
    if !ok {
        return fmt.Errorf("invalid session data type for phone %s: %T", phoneNumber, sessionInterface)
    }
    
    // 更新会话
    session.ChatIDs = chatIDs
    session.LastActive = time.Now()
    s.sessions.Store(phoneNumber, session)
    
    return nil
}

// DeleteSession 删除会话
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param phoneNumber 电话号码
// @return error 错误信息
func (s *sessionServiceImpl) DeleteSession(phoneNumber string) error {
    if phoneNumber == "" {
        return fmt.Errorf("phone number cannot be empty")
    }
    
    // 检查是否存在
    if _, exists := s.sessions.Load(phoneNumber); !exists {
        return fmt.Errorf("no session found for phone %s", phoneNumber)
    }
    
    // 删除会话
    s.sessions.Delete(phoneNumber)
    return nil
}

// ListActiveSessions 列出所有活跃会话
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @return []*Session 活跃会话列表
func (s *sessionServiceImpl) ListActiveSessions() []*Session {
    var activeSessions []*Session
    
    // 遍历所有会话
    s.sessions.Range(func(key, value interface{}) bool {
        session, ok := value.(*Session)
        if !ok {
            return true
        }
        
        // 检查会话是否活跃（最近30分钟内有活动）
        if time.Since(session.LastActive) < 30*time.Minute {
            activeSessions = append(activeSessions, session)
        }
        return true
    })
    
    return activeSessions
}

// GetConfiguredChatIDs 获取所有配置的群组ID
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @return []int64 群组ID列表
func (s *sessionServiceImpl) GetConfiguredChatIDs() []int64 {
    chatIDMap := make(map[int64]bool)
    
    // 遍历所有会话，收集唯一的群组ID
    s.sessions.Range(func(key, value interface{}) bool {
        session, ok := value.(*Session)
        if !ok {
            return true
        }
        
        for _, chatID := range session.ChatIDs {
            chatIDMap[chatID] = true
        }
        return true
    })
    
    // 转换为切片
    chatIDs := make([]int64, 0, len(chatIDMap))
    for chatID := range chatIDMap {
        chatIDs = append(chatIDs, chatID)
    }
    
    return chatIDs
}

/*
 * 关键算法说明：
 * 1. 会话管理：使用sync.Map存储和管理多个用户会话
 * 2. 活跃会话过滤：根据最后活跃时间筛选活跃会话
 * 3. 群组ID去重：使用map收集唯一的群组ID
 * 
 * 待优化事项：
 * 1. 会话持久化：将会话数据持久化到数据库
 * 2. 会话过期：添加自动过期和清理机制
 * 3. 并发控制：优化高并发场景下的性能
 * 
 * 兼容性说明：
 * 1. 使用Go标准库实现，无外部依赖
 * 2. 支持多会话并发访问
 */