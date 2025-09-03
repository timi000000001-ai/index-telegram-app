/**
 * 文件功能: 全局消息（Toast）Store 与便捷 API，统一应用内的提示/错误展示
 * 主要导出: messageStore, info, success, error, remove
 * 修改历史:
 * - 2025-09-02 v1.0.0 新增全局消息 store 与便捷方法
 * @author Trae AI
 * @date 2025-09-02
 * @version 1.0.0
 * © 项目组/组织
 */

import { writable } from 'svelte/store';

/**
 * 消息结构体
 * @typedef {{ id: string; type: 'info'|'success'|'error'; text: string; ttl: number }} MessageItem
 */

/**
 * 全局消息存储（队列）
 * @type {import('svelte/store').Writable<MessageItem[]>}
 */
export const messageStore = writable([]);

/**
 * 生成唯一 ID
 * @returns {string}
 */
function uid() {
  return Math.random().toString(36).slice(2) + Date.now().toString(36);
}

/**
 * 推送一条消息
 * @param {'info'|'success'|'error'} type 消息类型
 * @param {string} text 文本内容
 * @param {number} [ttl=3000] 存活毫秒数
 * @returns {string} 消息ID
 */
export function push(type, text, ttl = 3000) {
  const id = uid();
  const item = { id, type, text, ttl };
  messageStore.update((arr) => [...arr, item]);
  // 到期自动移除
  setTimeout(() => remove(id), ttl);
  return id;
}

/**
 * 按 ID 移除消息
 * @param {string} id 消息 ID
 * @returns {void}
 */
export function remove(id) {
  messageStore.update((arr) => arr.filter((m) => m.id !== id));
}

/**
 * 信息提示
 * @param {string} text 文本
 * @param {number} [ttl] 存活时间
 */
export function info(text, ttl) { return push('info', text, ttl); }

/**
 * 成功提示
 * @param {string} text 文本
 * @param {number} [ttl] 存活时间
 */
export function success(text, ttl) { return push('success', text, ttl); }

/**
 * 错误提示
 * @param {string} text 文本
 * @param {number} [ttl] 存活时间
 */
export function error(text, ttl) { return push('error', text, ttl); }

/*
 * 关键算法说明：基于简单数组队列 + setTimeout 到期清理
 * 待优化事项：支持可关闭按钮、不同位置与动画；合并重复消息
 * 兼容性说明：现代浏览器 setTimeout 与数组操作；SSR 安全（仅渲染依赖）
 */