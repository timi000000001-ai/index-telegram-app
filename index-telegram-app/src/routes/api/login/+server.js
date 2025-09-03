/**
 * 文件功能: 登录 Mock 接口（POST /api/login）
 * 主要导出: POST
 * 修改历史:
 * - 2025-09-02 v1.0.0 新增登录 mock 接口
 * @author Trae AI
 * @date 2025-09-02
 * @version 1.0.0
 * © 项目组/组织
 */

import { json } from '@sveltejs/kit';

/**
 * POST /api/login
 * @param {import('@sveltejs/kit').RequestEvent} event
 */
export async function POST(event) {
  await new Promise((r) => setTimeout(r, 200));
  const body = await event.request.json().catch(() => ({}));
  const { username, password } = body;
  if (!username || !password) {
    return json({ code: 400, message: '用户名或密码为空' }, { status: 400 });
  }
  return json({
    code: 200,
    data: {
      token: 'mock-token-123',
      user: { id: 'u001', name: username, role: 'tester' }
    }
  });
}

/*
 * 关键算法说明：无
 * 待优化事项：模拟错误次数上限、验证码
 * 兼容性说明：标准 JSON 输出
 */