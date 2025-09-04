import { json } from '@sveltejs/kit';

/**
 * @typedef {Object} Bot
 * @property {number} id
 * @property {string} name
 * @property {string} status
 * @property {string} uptime
 * @property {number} messages
 * @property {string} responseTime
 * @property {number} errors
 * @property {string} lastActivity
 * @property {string} createdAt
 */

/** @type {Bot[]} */
const mockBots = [
    {
        id: 1,
        name: '搜索机器人 #1',
        status: 'online',
        uptime: '2天5小时',
        messages: 1500,
        responseTime: '120ms',
        errors: 2,
        lastActivity: '2分钟前',
        createdAt: '2025-01-20'
    },
    {
        id: 2,
        name: '搜索机器人 #2',
        status: 'online',
        uptime: '1天12小时',
        messages: 890,
        responseTime: '95ms',
        errors: 0,
        lastActivity: '1分钟前',
        createdAt: '2025-01-19'
    },
    {
        id: 3,
        name: '搜索机器人 #3',
        status: 'starting',
        uptime: '3小时',
        messages: 45,
        responseTime: '150ms',
        errors: 1,
        lastActivity: '5分钟前',
        createdAt: '2025-01-22'
    },
    {
        id: 4,
        name: '索引机器人 #A',
        status: 'offline',
        uptime: 'N/A',
        messages: 0,
        responseTime: 'N/A',
        errors: 10,
        lastActivity: '3天前',
        createdAt: '2025-01-10'
    }
];

export async function GET() {
    await new Promise(resolve => setTimeout(resolve, 200)); // 模拟网络延迟
    return json({ success: true, bots: mockBots });
}