import { json } from '@sveltejs/kit';

/**
 * @type {import('./$types').RequestHandler}
 */
export async function POST({ request }) {
  const { code } = await request.json();

  // 在这里添加验证码的验证逻辑
  // 例如，检查 code 是否有效
  if (code === '123456') { // 这是一个示例，您需要替换为实际的验证逻辑
    // 验证成功，返回成功状态
    return json({ status: 'success' });
  } else {
    // 验证失败，返回错误信息
    return json({ status: 'error', error: '无效的验证码' }, { status: 401 });
  }
}