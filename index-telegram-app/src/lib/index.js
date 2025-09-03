// place files you want to import through the `$lib` alias in this folder.

/**
 * 模块功能: 提供统一的 API 请求封装与环境变量读取
 * 主要接口: apiFetch, API_BASE, AUTH_TOKEN
 * 修改历史:
 * - 2025-09-02 v1.0.0 新增 API 封装，接入环境变量 VITE_API_BASE / VITE_AUTH_TOKEN
 * - 2025-09-02 v1.1.0 默认使用同源地址，便于本地 mock 接口联调
 * @author Trae AI
 * @date 2025-09-02
 * @version 1.1.0
 * © 公司/组织
 */

/** 基础地址，优先使用环境变量 VITE_API_BASE；否则在浏览器默认同源，SSR 回退到 http://localhost:5174 */
export const API_BASE = (import.meta.env?.VITE_API_BASE)
  || (typeof window !== 'undefined' ? window.location.origin : 'http://localhost:5174');
/** 鉴权令牌（可选），来自环境变量 VITE_AUTH_TOKEN */
export const AUTH_TOKEN = import.meta.env?.VITE_AUTH_TOKEN || '';

/**
 * 统一 API 请求封装
 * @param {string} path 请求路径（相对 API_BASE），例如 "/api/search"
 * @param {{ method?: string, headers?: Record<string, string>, body?: any, query?: Record<string, string|number|boolean|null|undefined> }} [options]
 * @returns {Promise<{ ok: boolean, status: number, data?: any, error?: string }>}
 */
export async function apiFetch(path, options = {}) {
  const { method = 'GET', headers = {}, body, query } = options;
  const base = (API_BASE?.replace(/\/$/, '')) || '';
  // 当 path 是以 http(s) 开头时，直接使用完整地址；否则拼接 base
  const isAbsolute = /^https?:\/\//i.test(path);
  const url = isAbsolute ? new URL(path) : new URL(path, base + '/');
  if (query && typeof query === 'object') {
    for (const [k, v] of Object.entries(query)) {
      if (v !== undefined && v !== null) url.searchParams.set(k, String(v));
    }
  }
  const h = {
    'Content-Type': 'application/json',
    ...(AUTH_TOKEN ? { Authorization: `Bearer ${AUTH_TOKEN}` } : {}),
    ...headers
  };
  try {
    const resp = await fetch(url.toString(), {
      method,
      headers: h,
      body: body !== undefined && method !== 'GET' ? JSON.stringify(body) : undefined
    });
    const status = resp.status;
    const data = await (async () => { try { return await resp.json(); } catch { return null; } })();
    // 后端返回 code === 200 视为成功，否则按错误处理
    const backendOk = data && typeof data === 'object' && ('code' in data ? data.code === 200 : true);
    if (!resp.ok || !backendOk) {
      const error = (data && data.message) || resp.statusText || 'Request Error';
      return { ok: false, status, data, error };
    }
    return { ok: true, status, data };
  } catch (e) {
    const msg = e && typeof e === 'object' && 'message' in e ? e.message : String(e);
    return { ok: false, status: 0, error: String(msg) };
  }
}