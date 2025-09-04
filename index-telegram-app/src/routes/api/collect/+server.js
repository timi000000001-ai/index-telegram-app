/**
 * 文件功能: 采集配置 Mock 接口（POST /api/collect）
 * 主要导出: POST
 * 修改历史:
 * - 2025-09-05 v1.0.0 新增采集配置 mock
 * @author Trae AI
 * @date 2025-09-05
 * @version 1.0.0
 * © 项目组/组织
 */

import { json } from '@sveltejs/kit';

/**
 * POST /api/collect
 * @param {import('@sveltejs/kit').RequestEvent} event
 */
export async function POST(event) {
  const { chat_ids } = await event.request.json();

  await new Promise((r) => setTimeout(r, 150));

  if (!chat_ids || !Array.isArray(chat_ids) || chat_ids.length === 0) {
    return json({ status: 'error', message: '群组 ID 列表不能为空' }, { status: 400 });
  }

  console.log(`Mock: 收到配置请求，群组ID: ${chat_ids.join(', ')}`);

  return json({
    status: 'success',
    message: `已成功配置 ${chat_ids.length} 个群组。`,
  });
}

/*
 * 关键算法说明：模拟配置提交流程
 * 待优化事项：无
 * 兼容性说明：标准 JSON 输出，已与 Postman 定义对齐
 */