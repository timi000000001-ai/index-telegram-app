/*
 * 文件功能描述：搜索服务，处理搜索请求和结果处理
 * 主要类/接口说明：SearchService接口及其实现
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
    "net/url"
    "strconv"

    "github.com/go-resty/resty/v2"
)

// SearchResult 搜索结果结构
// @author fcj
// @date 2023-11-15
// @version 1.0.0
type SearchResult struct {
    Hits       []map[string]interface{} `json:"hits"`
    TotalHits  int                      `json:"totalHits"`
    TotalPages int                      `json:"totalPages"`
}

// SearchConfig 搜索服务配置
// @author fcj
// @date 2023-11-15
// @version 1.0.0
type SearchConfig struct {
    MeilisearchURL string
    MeilisearchKey string
}

// SearchService 定义搜索服务接口
// @author fcj
// @date 2023-11-15
// @version 1.0.0
type SearchService interface {
    // Search 执行搜索操作
    Search(query string, page string, limit string, filter string) (*SearchResult, error)
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

// Search 执行搜索操作
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param query 搜索查询
// @param page 页码
// @param limit 每页结果数
// @param filter 过滤条件
// @return *SearchResult 搜索结果
// @return error 错误信息
func (s *searchServiceImpl) Search(query string, page string, limit string, filter string) (*SearchResult, error) {
    // 设置默认值
    if page == "" {
        page = "1"
    }
    if limit == "" {
        limit = "5"
    }

    // 构建查询参数
    params := url.Values{}
    params.Set("q", query)
    params.Set("offset", fmt.Sprintf("%d", (atoi(page)-1)*atoi(limit)))
    params.Set("limit", limit)

    // 根据过滤条件设置筛选参数
    switch filter {
    case "group":
        params.Set("filter", "chat_type = 'group'")
    case "channel":
        params.Set("filter", "chat_type = 'channel'")
    case "bot":
        params.Set("filter", "sender_id LIKE '%bot%'")
    case "message", "":
        // 默认不添加筛选
    }

    // 执行搜索请求
    resp, err := s.client.R().
        SetHeader("Authorization", "Bearer "+s.config.MeilisearchKey).
        SetQueryParamsFromValues(params).
        Get(s.config.MeilisearchURL + "/indexes/messages/search")
    if err != nil {
        return nil, fmt.Errorf("search failed: %w", err)
    }

    // 解析搜索结果
    var result struct {
        Hits      []map[string]interface{} `json:"hits"`
        TotalHits int                      `json:"totalHits"`
    }
    if err := json.Unmarshal(resp.Body(), &result); err != nil {
        return nil, fmt.Errorf("failed to parse search result: %w", err)
    }

    // 计算总页数
    totalPages := (result.TotalHits + atoi(limit) - 1) / atoi(limit)

    // 构建搜索结果
    searchResult := &SearchResult{
        Hits:       result.Hits,
        TotalHits:  result.TotalHits,
        TotalPages: totalPages,
    }

    return searchResult, nil
}

// atoi 将字符串转换为整数
// @author fcj
// @date 2023-11-15
// @version 1.0.0
// @param s 要转换的字符串
// @return int 转换后的整数
func atoi(s string) int {
    i, _ := strconv.Atoi(s)
    return i
}

/*
 * 关键算法说明：
 * 1. 搜索分页：使用offset和limit实现分页
 * 2. 过滤条件：根据不同的过滤条件构建Meilisearch查询
 * 
 * 待优化事项：
 * 1. 错误处理：改进错误处理和恢复机制
 * 2. 缓存机制：实现搜索结果缓存
 * 3. 高级搜索：支持更复杂的搜索条件
 * 
 * 兼容性说明：
 * 1. 依赖Meilisearch API
 * 2. 需要有效的Meilisearch认证令牌
 */