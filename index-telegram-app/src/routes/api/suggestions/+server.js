/**
 * 文件功能: 搜索建议 Mock 接口（GET /api/suggestions?q=xxx）
 * 主要导出: GET
 * 修改历史:
 * - 2025-09-02 v1.0.0 新增搜索建议 mock
 * - 2025-09-02 v1.0.1 修正类型注释，避免 linter 报错
 * @author Trae AI
 * @date 2025-09-02
 * @version 1.0.1
 * © 项目组/组织
 */

import { json } from '@sveltejs/kit';

/**
 * GET /api/suggestions?q=xxx
 * @param {import('@sveltejs/kit').RequestEvent} event
 */
export async function GET(event) {
  const q = event.url.searchParams.get('q') || '';
  await new Promise((r) => setTimeout(r, 120));
  const base = ['news', 'sports', 'music', 'tech', 'movie', 'Svelte', 'Kit', 'Telegram', 'AI'];
  const suggestions = base.filter((w) => w.toLowerCase().includes(q.toLowerCase())).slice(0, 8);
  return json({ code: 200, data: { suggestions } });
}

/*
 * 关键算法说明：简单匹配
 * 待优化事项：根据热度排序
 * 兼容性说明：标准 JSON 输出
 */