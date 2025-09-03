/**
 * 文件功能: 机器人状态 Mock 接口（GET /api/bot/status）
 * 主要导出: GET
 * 修改历史:
 * - 2025-09-02 v1.0.0 新增机器人状态 mock
 * @author Trae AI
 * @date 2025-09-02
 * @version 1.0.0
 * © 项目组/组织
 */

import { json } from '@sveltejs/kit';

/**
 * GET /api/bot/status
 * @returns {Promise<Response>} 响应包含 bots 列表
 */
export async function GET() {
  // 模拟网络延迟
  await new Promise((r) => setTimeout(r, 200));
  return json({
    code: 200,
    data: {
      bots: [
        { id: 'bot-1', name: 'tg_search_bot', status: 'online', lastActive: new Date().toISOString() }
      ]
    }
  });
}

/*
 * 关键算法说明：无
 * 待优化事项：可根据环境变量切换不同状态
 * 兼容性说明：标准 JSON 输出
 */