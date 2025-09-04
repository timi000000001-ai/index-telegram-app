/**
 * 文件功能: 热门搜索 Mock 接口 (GET /api/search/trending)
 * 主要导出: GET
 * @author Trae AI
 * @date 2025-09-06
 * @version 1.0.0
 */

import { json } from '@sveltejs/kit';

/**
 * GET /api/search/trending
 * @param {import('@sveltejs/kit').RequestEvent} event
 */
export async function GET(event) {
  await new Promise((r) => setTimeout(r, 50)); // 模拟网络延迟
  const trendingData = [
    { keyword: 'SvelteKit', trend: 'hot', count: 1200, rank: 1, category: 'tech' },
    { keyword: 'Telegram API', trend: 'up', count: 980, rank: 2, category: 'dev' },
    { keyword: 'AI 助手', trend: 'up', count: 850, rank: 3, category: 'ai' },
    { keyword: '数据可视化', trend: 'stable', count: 720, rank: 4, category: 'tech' },
    { keyword: 'Web3', trend: 'down', count: 650, rank: 5, category: 'blockchain' },
    { keyword: '机器学习', trend: 'hot', count: 580, rank: 6, category: 'ai' },
    { keyword: '区块链', trend: 'stable', count: 520, rank: 7, category: 'blockchain' },
    { keyword: '前端开发', trend: 'up', count: 480, rank: 8, category: 'dev' }
  ];
  return json(trendingData);
}