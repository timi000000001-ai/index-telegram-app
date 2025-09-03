/**
 * 文件功能：提供机器人状态的 API 接口（GET /api/bots/status），用于前端展示“我的机器人状态”模块。
 * 主要导出：GET 请求处理函数
 * 修改历史：
 * - 2025-09-02 v1.0.1 修正返回类型注释
 * - 2025-09-02 v1.0.0 初始化实现，返回模拟数据，修复前端 404 报错
 *
 * 元信息
 * @author AI Assistant
 * @date 2025-09-02
 * @version 1.0.1
 * © 项目组/组织
 */

/**
 * GET /api/bots/status
 * 参数：无
 * 返回：{ success: boolean, bots: Array<BotStatus>, error?: string }
 * 说明：返回机器人详细状态信息，包含运行时间、消息处理量等统计数据
 * @returns {Promise<Response>}
 */
export async function GET() {
  try {
    const data = {
      success: true,
      bots: [
        {
          id: 1,
          name: '搜索机器人 #1',
          status: 'online',
          uptime: '2天5小时',
          messages: 1500,
          responseTime: '120ms',
          errors: 2,
          lastActivity: '2分钟前',
          createdAt: '2025-01-20',
          description: '主要搜索服务机器人'
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
          createdAt: '2025-01-19',
          description: '备用搜索服务机器人'
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
          createdAt: '2025-01-22',
          description: '新部署的测试机器人'
        }
      ]
    };
    
    return new Response(JSON.stringify(data), {
      headers: { 'content-type': 'application/json; charset=utf-8' }
    });
  } catch (error) {
    return new Response(JSON.stringify({
      success: false,
      error: '获取机器人状态失败',
      bots: []
    }), {
      status: 500,
      headers: { 'content-type': 'application/json; charset=utf-8' }
    });
  }
}