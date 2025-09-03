/**
 * 文件功能: 搜索结果 Mock 接口（GET /api/search）
 * 主要导出: GET
 * 修改历史:
 * - 2025-09-02 v1.0.0 新增搜索结果 mock（分页/排序/过滤）
 * @author Trae AI
 * @date 2025-09-02
 * @version 1.0.0
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
  const size = parseInt(url.searchParams.get('size') || '20', 10);
  const sort = url.searchParams.get('sort') || 'relevance';
  await new Promise((r) => setTimeout(r, 180));

  // 生成模拟数据
  const total = 137;
  const start = (page - 1) * size;
  const end = Math.min(start + size, total);
  const results = Array.from({ length: Math.max(0, end - start) }).map((_, i) => {
    const id = `${start + i + 1}`;
    const type = ['text', 'media', 'link'][Math.floor(Math.random() * 3)];
    return {
      id,
      type,
      source: ['group:科技圈', 'channel:日报', 'user:Alice'][Math.floor(Math.random() * 3)],
      timestamp: new Date(Date.now() - Math.random() * 7 * 24 * 3600 * 1000).toISOString(),
      relevance: Math.round(Math.random() * 100) / 100,
      content: `【${q || '关键词'}】的模拟结果 ${id} —— 这是一条用于联调的占位内容。`
    };
  });

  // 简单排序模拟
  if (sort === 'time') results.sort((a, b) => +new Date(b.timestamp) - +new Date(a.timestamp));
  if (sort === 'heat') results.sort((a, b) => (b.relevance ?? 0) - (a.relevance ?? 0));

  return json({ code: 200, data: { total, page, size, results } });
}

/*
 * 关键算法说明：模拟分页与排序
 * 待优化事项：根据过滤参数生成更真实数据
 * 兼容性说明：标准 JSON 输出
 */