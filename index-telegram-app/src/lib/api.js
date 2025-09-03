/**
 * API 接口封装模块
 * @author 前端工程师
 * @date 2024-01-23
 * @version 1.0.0
 * © Telegram Search App
 */

/**
 * 统一的 API 请求封装函数
 * @param {string} url - 请求URL
 * @param {RequestInit} options - 请求选项
 * @returns {Promise<any>} 返回响应数据
 */
export async function apiFetch(url, options = {}) {
    try {
        const response = await fetch(url, {
            headers: {
                'Content-Type': 'application/json',
                ...(options.headers || {})
            },
            ...options
        });
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        return await response.json();
    } catch (error) {
        console.error('API请求失败:', error);
        throw error;
    }
}

/**
 * 搜索相关API
 */
export const searchAPI = {
    /**
     * 执行搜索
     * @param {Record<string, string>} params - 搜索参数
     * @returns {Promise<Object>} 搜索结果
     */
    async search(params) {
        const queryString = new URLSearchParams(params).toString();
        return apiFetch(`/api/search?${queryString}`);
    },
    
    /**
     * 获取搜索建议
     * @param {string} query - 查询关键词
     * @returns {Promise<string[]>} 建议列表
     */
    async getSuggestions(query) {
        return apiFetch(`/api/suggestions?q=${encodeURIComponent(query)}`);
    },
    
    /**
     * 获取热门关键词
     * @returns {Promise<string[]>} 热门关键词列表
     */
    async getTrending() {
        return apiFetch('/api/trending');
    }
};

/**
 * 机器人相关API
 */
export const botAPI = {
    /**
     * 获取机器人列表
     * @returns {Promise<any[]>} 机器人列表
     */
    async getBots() {
        return apiFetch('/api/bots');
    },
    
    /**
     * 获取机器人详情
     * @param {string} botId - 机器人ID
     * @returns {Promise<Object>} 机器人详情
     */
    async getBotDetail(botId) {
        return apiFetch(`/api/bots/${botId}`);
    }
};

/**
 * 用户相关API
 */
export const userAPI = {
    /**
     * 用户登录
     * @param {Object} credentials - 登录凭据
     * @returns {Promise<Object>} 登录结果
     */
    async login(credentials) {
        return apiFetch('/api/auth/login', {
            method: 'POST',
            body: JSON.stringify(credentials)
        });
    },
    
    /**
     * 获取用户信息
     * @returns {Promise<Object>} 用户信息
     */
    async getProfile() {
        return apiFetch('/api/user/profile');
    }
};

// 默认导出
export default {
    apiFetch,
    searchAPI,
    botAPI,
    userAPI
};

/**
 * 关键算法说明：
 * - 使用统一的错误处理机制
 * - 支持请求头自定义
 * - 自动处理JSON响应
 * 
 * 待优化事项：
 * - 添加请求拦截器
 * - 实现请求重试机制
 * - 添加缓存策略
 * 
 * 兼容性说明：
 * - 支持现代浏览器的fetch API
 * - 需要polyfill支持旧版浏览器
 */