/**
 * 文件功能: 用户信息 Mock 接口 (GET /api/user/profile)
 * 主要导出: GET
 * @author Trae AI
 * @date 2025-09-06
 * @version 1.0.0
 */

import { json } from '@sveltejs/kit';

/**
 * GET /api/user/profile
 * @param {import('@sveltejs/kit').RequestEvent} event
 */
export async function GET(event) {
  await new Promise((r) => setTimeout(r, 80)); // 模拟网络延迟
  const userProfile = {
    id: 'u007',
    name: 'Trae AI',
    email: 'trae@example.com',
    avatar: 'https://avatars.githubusercontent.com/u/106141222?s=200&v=4',
  };
  return json(userProfile);
}