/*
 * 文件功能描述：搜索服务，处理数据检索和索引功能
 * 主要类/接口说明：SearchService接口及其实现
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
    "net/url"

    "github.com/go-resty/resty/v2"
)

// SearchService 定义搜索服务接口
// @author fcj
// @date 2023-11-15
// @version 1.0.0
type SearchService interface {
    // IndexToMeilisearch 将数据索引到Meilisearch
    IndexToMeilisearch(data map[string]interface{}) error
    
    // Search 搜索数据
    Search(query string, page int, limit int, filter string) ([]byte, error)
}

// SearchConfig 搜索服务配置
// @author fcj
// @date 2023-11-15
// @version 1.0.0
type SearchConfig struct {
    MeilisearchURL  string
    MeilisearchKey string
    ManagementServiceURL string
}

// searchServiceImpl 实现SearchService接口
// @author fcj
// @date 2023-11-15
// @version 1.0.0
type searchServiceImpl struct {
    config SearchConfig
    client *resty.Client
}

// NewSearchService 创建新的搜索服务实例
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param config 搜索服务配置
// @return SearchService 搜索服务实例
func NewSearchService(config SearchConfig) SearchService {
    return &searchServiceImpl{
        config: config,
        client: resty.New(),
    }
}

// IndexToMeilisearch 将数据索引到Meilisearch
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param data 需要索引的数据
// @return error 错误信息
func (s *searchServiceImpl) IndexToMeilisearch(data map[string]interface{}) error {
    resp, err := s.client.R().
        SetHeader("Authorization", "Bearer "+s.config.MeilisearchKey).
        SetBody([]map[string]interface{}{data}).
        Post(s.config.MeilisearchURL + "/indexes/messages/documents")
    if err != nil {
        log.Printf("Failed to connect to Meilisearch: %v", err)
        return fmt.Errorf("failed to connect to Meilisearch: %w", err)
    }
    if resp.StatusCode() != 202 {
        log.Printf("Meilisearch returned non-202 status: %d", resp.StatusCode())
        return fmt.Errorf("Meilisearch returned status code: %d", resp.StatusCode())
    }
    return nil
}

// Search 搜索数据
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param query 搜索查询
// @param page 页码
// @param limit 每页结果数
// @param filter 过滤条件
// @return []byte 搜索结果
// @return error 错误信息
func (s *searchServiceImpl) Search(query string, page int, limit int, filter string) ([]byte, error) {
    // 构建查询参数
    params := url.Values{}
    params.Set("q", query)
    params.Set("page", fmt.Sprintf("%d", page))
    params.Set("limit", fmt.Sprintf("%d", limit))
    if filter != "" {
        params.Set("filter", filter)
    }
    
    // 调用管理服务的搜索API
    resp, err := s.client.R().
        SetHeader("Authorization", "Bearer YOUR_AUTH_TOKEN").
        SetQueryParamsFromValues(params).
        Get(s.config.ManagementServiceURL + "/api/search")
    if err != nil {
        log.Printf("Failed to search: %v", err)
        return nil, fmt.Errorf("failed to search: %w", err)
    }
    
    return resp.Body(), nil
}

/*
 * 关键算法说明：
 * 1. 索引操作：将数据索引到Meilisearch搜索引擎
 * 2. 搜索操作：通过管理服务API进行搜索
 * 
 * 待优化事项：
 * 1. 缓存机制：实现搜索结果缓存
 * 2. 批量索引：支持批量数据索引
 * 3. 高级搜索：支持更复杂的搜索条件
 * 
 * 兼容性说明：
 * 1. 依赖Meilisearch API和管理服务API
 * 2. 需要有效的认证令牌
 */