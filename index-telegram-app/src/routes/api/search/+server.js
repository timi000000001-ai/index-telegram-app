/**
 * 文件功能: 搜索结果 Mock 接口（GET /api/search）
 * 主要导出: GET
 * 修改历史:
 * - 2025-09-02 v1.0.0 新增搜索结果 mock（分页/排序/过滤）
 * - 2025-09-05 v1.1.0 修正接口参数与响应结构，与 Postman 定义对齐
 * @author Trae AI
 * @date 2025-09-05
 * @version 1.1.0
 * © 项目组/组织
 */

import { json } from '@sveltejs/kit';

/**
 * GET /api/search
 * @param {import('@sveltejs/kit').RequestEvent} event
 */
export async function GET(event) {
  const url = event.url;
  const q = url.searchParams.get('q') || '';
  const page = parseInt(url.searchParams.get('page') || '1', 10);
  const limit = parseInt(url.searchParams.get('limit') || '20', 10);
  const filter = url.searchParams.get('filter') || 'all';
  const sort = url.searchParams.get('sort') || 'relevance';
  await new Promise((r) => setTimeout(r, 180));

  // 生成模拟数据
  const total = 137;
  const pages = Math.ceil(total / limit);
  const start = (page - 1) * limit;
  const end = Math.min(start + limit, total);
  const results = Array.from({ length: Math.max(0, end - start) }).map((_, i) => {
    const id = `${start + i + 1}`;
    const chat_type = ['group', 'channel', 'bot'][Math.floor(Math.random() * 3)];
    return {
      id: `rec${id}`,
      content: `【${q || '关键词'}】的模拟结果 ${id} —— 这是一条用于联调的占位内容。`,
      sender: `user_${id}`,
      source: `示例${chat_type}`,
      type: filter === 'all' ? chat_type : filter,
      timestamp: new Date(Date.now() - Math.random() * 7 * 24 * 3600 * 1000).toISOString(),
      relevance: Math.round(Math.random() * 100) / 100, // 用于排序
    };
  });

  // 简单排序模拟
  if (sort === 'date') results.sort((a, b) => +new Date(b.timestamp) - +new Date(a.timestamp));
  else results.sort((a, b) => b.relevance - a.relevance);

  return json({ results, total, page, limit, pages });
}

/*
 * 关键算法说明：模拟分页与排序
 * 待优化事项：无
 * 兼容性说明：标准 JSON 输出，已与 Postman 定义对齐
 */