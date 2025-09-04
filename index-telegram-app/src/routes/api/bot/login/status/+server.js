import { json } from '@sveltejs/kit';
import { loginCodes } from '../../../../../lib/botLoginStore.js';

/** @type {import('./$types').RequestHandler} */
export async function GET({ url }) {
  const token = url.searchParams.get('token');

  if (!token) {
    return json({ error: 'Token is required' }, { status: 400 });
  }

  // Find the code associated with the token
  let foundCode = null;
  for (const [code, value] of loginCodes.entries()) {
    if (value.token === token) {
      foundCode = code;
      break;
    }
  }

  if (foundCode) {
    const loginInfo = loginCodes.get(foundCode);
    return json({ status: loginInfo.status });
  } else {
    return json({ status: 'expired' }, { status: 404 });
  }
}