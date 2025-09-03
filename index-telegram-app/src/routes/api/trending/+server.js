/**
 * 文件功能: 热门关键字 Mock 接口（GET /api/trending）
 * 主要导出: GET
 * 修改历史:
 * - 2025-09-02 v1.0.0 新增热门关键字 mock
 * @author Trae AI
 * @date 2025-09-02
 * @version 1.0.0
 * © 项目组/组织
 */

import { json } from '@sveltejs/kit';

export async function GET() {
  await new Promise((r) => setTimeout(r, 150));
  return json({
    code: 200,
    data: {
      trending: [
        { keyword: 'SvelteKit', count: 120, trend: 'up' },
        { keyword: 'Telegram', count: 96, trend: 'flat' },
        { keyword: 'AI', count: 80, trend: 'down' }
      ]
    }
  });
}

/*
 * 关键算法说明：无
 * 待优化事项：支持分页与时间窗口参数
 * 兼容性说明：标准 JSON 输出
 */