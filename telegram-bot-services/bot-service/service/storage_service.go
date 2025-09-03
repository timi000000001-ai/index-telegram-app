/*
 * 文件功能描述：存储服务，处理数据存储和索引功能
 * 主要类/接口说明：StorageService接口及其实现
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

    "github.com/go-resty/resty/v2"
)

// StorageService 定义存储服务接口
// @author fcj
// @date 2023-11-15
// @version 1.0.0
type StorageService interface {
    // SaveToPocketBase 保存消息到PocketBase
    SaveToPocketBase(data map[string]interface{}) error
    
    // IndexToMeilisearch 将消息索引到Meilisearch
    IndexToMeilisearch(data map[string]interface{}) error
    
    // SaveAndIndex 保存并索引消息
    SaveAndIndex(data map[string]interface{}) error
}

// StorageConfig 存储服务配置
// @author fcj
// @date 2023-11-15
// @version 1.0.0
type StorageConfig struct {
    PocketBaseURL    string
    MeilisearchURL   string
    MeilisearchToken string
}

// storageServiceImpl 实现StorageService接口
// @author fcj
// @date 2023-11-15
// @version 1.0.0
type storageServiceImpl struct {
    config StorageConfig
    client *resty.Client
}

// NewStorageService 创建新的存储服务实例
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param config 存储服务配置
// @return StorageService 存储服务实例
func NewStorageService(config StorageConfig) StorageService {
    return &storageServiceImpl{
        config: config,
        client: resty.New(),
    }
}

// SaveToPocketBase 保存消息到PocketBase
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param data 需要保存的消息数据
// @return error 错误信息
func (s *storageServiceImpl) SaveToPocketBase(data map[string]interface{}) error {
    resp, err := s.client.R().
        SetBody(data).
        Post(s.config.PocketBaseURL + "/api/collections/messages/records")
    if err != nil {
        log.Printf("Failed to save to PocketBase: %v", err)
        return fmt.Errorf("failed to connect to PocketBase: %w", err)
    }
    if resp.StatusCode() != 200 && resp.StatusCode() != 201 {
        log.Printf("PocketBase returned non-200/201 status: %d, body: %s", resp.StatusCode(), resp.Body())
        return fmt.Errorf("PocketBase returned status code: %d", resp.StatusCode())
    }
    log.Printf("Successfully saved message to PocketBase")
    return nil
}

// IndexToMeilisearch 将消息索引到Meilisearch
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param data 需要索引的消息数据
// @return error 错误信息
func (s *storageServiceImpl) IndexToMeilisearch(data map[string]interface{}) error {
    resp, err := s.client.R().
        SetHeader("Authorization", "Bearer "+s.config.MeilisearchToken).
        SetBody([]map[string]interface{}{data}).
        Post(s.config.MeilisearchURL + "/indexes/messages/documents")
    if err != nil {
        log.Printf("Failed to connect to Meilisearch: %v", err)
        return fmt.Errorf("failed to connect to Meilisearch: %w", err)
    }
    if resp.StatusCode() != 202 {
        log.Printf("Meilisearch returned non-202 status: %d, body: %s", resp.StatusCode(), resp.Body())
        return fmt.Errorf("Meilisearch returned status code: %d", resp.StatusCode())
    }
    log.Printf("Successfully indexed message to Meilisearch")
    return nil
}

// SaveAndIndex 保存并索引消息
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param data 需要保存和索引的消息数据
// @return error 错误信息
func (s *storageServiceImpl) SaveAndIndex(data map[string]interface{}) error {
    // 保存到PocketBase
    if err := s.SaveToPocketBase(data); err != nil {
        return fmt.Errorf("failed to save to PocketBase: %w", err)
    }
    
    // 索引到Meilisearch
    if err := s.IndexToMeilisearch(data); err != nil {
        log.Printf("Warning: Failed to index to Meilisearch: %v", err)
        // 不返回错误，因为索引失败不应阻止整个流程
    }
    
    return nil
}

/*
 * 关键算法说明：
 * 1. 数据存储：使用REST API将数据保存到PocketBase
 * 2. 数据索引：使用REST API将数据索引到Meilisearch
 * 
 * 待优化事项：
 * 1. 批量操作：支持批量保存和索引
 * 2. 错误重试：添加重试机制处理临时网络故障
 * 3. 连接池：优化HTTP客户端连接池
 * 
 * 兼容性说明：
 * 1. 依赖PocketBase和Meilisearch API
 * 2. 需要有效的服务地址和认证令牌
 */