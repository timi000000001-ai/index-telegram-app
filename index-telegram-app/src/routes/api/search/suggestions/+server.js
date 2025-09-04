/**
 * 文件功能: 搜索建议 Mock 接口 (GET /api/search/suggestions)
 * 主要导出: GET
 * @author Trae AI
 * @date 2025-09-06
 * @version 1.0.0
 */

import { json } from '@sveltejs/kit';

/**
 * GET /api/search/suggestions
 * @param {import('@sveltejs/kit').RequestEvent} event
 */
export async function GET(event) {
  const url = event.url;
  const q = url.searchParams.get('q') || '';
  await new Promise((r) => setTimeout(r, 20)); // 模拟网络延迟
  const suggestions = [
    `${q} 教程`,
    `${q} 案例`,
    `${q} 新闻`,
    `${q} 最佳实践`,
    `如何学习 ${q}`
  ];
  return json(suggestions);
}