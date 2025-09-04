/**
 * 文件功能: 登录 Mock 接口（POST /api/login）
 * 主要导出: POST
 * 修改历史:
 * - 2025-09-02 v1.0.0 新增登录 mock
 * - 2025-09-05 v1.1.0 修正为两步登录，与 Postman 定义对齐
 * @author Trae AI
 * @date 2025-09-05
 * @version 1.1.0
 * © 项目组/组织
 */

import { json } from '@sveltejs/kit';

/**
 * POST /api/login
 * @param {import('@sveltejs/kit').RequestEvent} event
 */
export async function POST(event) {
  const { phone_number, code } = await event.request.json();

  await new Promise((r) => setTimeout(r, 260));

  // 如果没有 code，模拟发送验证码
  if (!code) {
    if (!phone_number) {
      return json({ status: 'error', message: '电话号码不能为空' }, { status: 400 });
    }
    return json({
      status: 'success',
      message: `验证码已发送到 ${phone_number}，请查收。`,
    });
  }

  // 如果有 code，模拟验证
  if (code === '12345') {
    return json({
      status: 'success',
      message: '登录成功',
      token: `mock-token-${Date.now()}`,
      user: { name: 'Mock User', phone: phone_number },
    });
  } else {
    return json({ status: 'error', message: '验证码错误' }, { status: 401 });
  }
}

/*
 * 关键算法说明：模拟两步登录流程
 * 待优化事项：无
 * 兼容性说明：标准 JSON 输出，已与 Postman 定义对齐
 */