import { json } from '@sveltejs/kit';
import { loginCodes } from '$lib/botLoginStore.js';

export async function POST() {
  const code = Math.floor(100000 + Math.random() * 900000).toString();
  const token = `token-${Date.now()}`;
  loginCodes.set(code, { token, status: 'pending' });

  // In a real app, you might want to associate this with a user session
  // and set an expiration for the code.

  return json({ loginCode: code, token });
}