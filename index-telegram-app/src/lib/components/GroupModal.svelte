<!--
 * @file GroupModal.svelte
 * @description 群组详情弹窗组件 - 显示群组/频道的统计信息、活跃用户和活跃度分布
 * @author fcj
 * @date 2025-01-22
 * @version 1.0.0
 * @copyright © 2025 Telegram Search Platform
-->

<script>
  /**
   * 组件属性
   * @param {boolean} show - 控制弹窗显示/隐藏
   * @param {Object|null} groupData - 群组数据对象
   */
  export let show = false;
  export let groupData = {
    name: '',
    username: '',
    description: '',
    createdAt: '',
    stats: {
      totalMembers: 0,
      onlineMembers: 0,
      messages24h: 0,
      avgDaily: 0,
      activityRate: 0,
      groupScore: 0
    },
    topUsers: [
      { username: '', messages: 0, activity: 0 }
    ],
    hourlyActivity: [0],
    lastUpdate: ''
  };

  /**
   * 事件派发器
   */
  import { createEventDispatcher } from 'svelte';
  const dispatch = createEventDispatcher();

  /**
   * 关闭弹窗处理函数
   * 派发close事件给父组件
   */
  function closeModal() {
    dispatch('close');
  }

  /**
   * 处理遮罩层点击事件
   * @param {Event} event - 点击事件对象
   */
  function handleOverlayClick(event) {
    // 检查点击是否在弹窗内容区域外
    if (event.currentTarget && event.target) {
      const modalContent = /** @type {HTMLElement} */ (event.currentTarget).querySelector('[role="document"]');
      if (modalContent && !modalContent.contains(/** @type {Node} */ (event.target))) {
        closeModal();
      }
    }
  }

  /**
   * 处理键盘ESC键关闭弹窗
   * @param {KeyboardEvent} event - 键盘事件对象
   */
  function handleKeydown(event) {
    if (event.key === 'Escape') {
      closeModal();
    }
  }
</script>

<!-- 键盘事件监听 -->
<svelte:window on:keydown={handleKeydown} />

<!-- 群组详情弹窗 -->
{#if show && groupData}
  <div class="fixed inset-0 bg-black/60 flex items-center justify-center z-50 backdrop-blur-sm p-5" role="dialog" aria-modal="true" tabindex="-1" on:click={handleOverlayClick} on:keydown={handleKeydown}>
    <div class="bg-white rounded-3xl max-w-4xl w-full max-h-[calc(100vh-40px)] overflow-hidden shadow-2xl animate-in slide-in-from-bottom-4 duration-300" role="document">
      <!-- 弹窗头部 -->
      <div class="flex items-center justify-between px-8 py-6 border-b border-gray-200 bg-gradient-to-r from-indigo-500 to-purple-600 text-white">
        <h2 class="text-2xl font-semibold m-0">群组动态统计详情</h2>
        <button class="bg-white/20 border-0 rounded-full w-10 h-10 flex items-center justify-center cursor-pointer transition-colors hover:bg-white/30" on:click={closeModal} title="关闭">
          <span class="text-2xl text-white leading-none">×</span>
        </button>
      </div>
      
      <!-- 弹窗内容 -->
      <div class="p-8 overflow-y-auto max-h-[calc(100vh-200px)]">
        <!-- 群组基本信息 -->
        <div class="flex items-center gap-5 mb-8 p-6 bg-gradient-to-br from-slate-50 to-gray-100 rounded-2xl">
          <div class="flex-shrink-0">
            <div class="w-20 h-20 rounded-full bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center shadow-lg shadow-indigo-500/30">
              <span class="text-white text-3xl">▶</span>
            </div>
          </div>
          <div class="flex-1">
            <h3 class="text-xl font-semibold text-gray-900 mb-1">{groupData.name}</h3>
            <p class="text-gray-600 text-sm mb-2">{groupData.username} · 创建于 {groupData.createdAt}</p>
            <p class="text-gray-700 text-sm leading-relaxed">{groupData.description}</p>
          </div>
        </div>
        
        <!-- 统计数据卡片 -->
        <div class="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
          <div class="bg-white border border-gray-200 rounded-xl p-4 text-center shadow-sm hover:shadow-md transition-shadow">
            <div class="text-2xl font-bold text-indigo-600 mb-1">{groupData.stats.totalMembers.toLocaleString()}</div>
            <div class="text-sm font-medium text-gray-700 mb-2">总成员数</div>
            <div class="text-lg font-semibold text-green-600">{groupData.stats.onlineMembers}</div>
            <div class="text-xs text-gray-500">在线成员</div>
          </div>
          <div class="bg-white border border-gray-200 rounded-xl p-4 text-center shadow-sm hover:shadow-md transition-shadow">
            <div class="text-2xl font-bold text-blue-600 mb-1">{groupData.stats.messages24h}</div>
            <div class="text-sm font-medium text-gray-700 mb-2">24h消息数</div>
            <div class="text-lg font-semibold text-green-600">{groupData.stats.avgDaily}</div>
            <div class="text-xs text-gray-500">活跃用户</div>
          </div>
          <div class="bg-white border border-gray-200 rounded-xl p-4 text-center shadow-sm hover:shadow-md transition-shadow">
            <div class="text-2xl font-bold text-orange-600 mb-1">{groupData.stats.avgDaily}.8</div>
            <div class="text-sm font-medium text-gray-700 mb-2">平均每日/用户</div>
            <div class="text-lg font-semibold text-green-600">{groupData.stats.activityRate}%</div>
            <div class="text-xs text-gray-500">活跃度</div>
          </div>
          <div class="bg-white border border-gray-200 rounded-xl p-4 text-center shadow-sm hover:shadow-md transition-shadow">
            <div class="text-2xl font-bold text-purple-600 mb-1">{groupData.stats.groupScore}</div>
            <div class="text-sm font-medium text-gray-700 mb-2">群组评分</div>
            <div class="text-lg font-semibold text-green-600">{groupData.stats.activityRate}%</div>
            <div class="text-xs text-gray-500">活跃度</div>
          </div>
        </div>
        
        <!-- 内容区域：活跃用户TOP5 和 24小时活跃度分布 -->
        <div class="grid lg:grid-cols-2 gap-6 mb-8">
          <!-- 活跃用户TOP5 -->
          <div class="bg-white border border-gray-200 rounded-xl p-6">
            <h4 class="text-lg font-semibold text-gray-900 mb-4">活跃用户TOP5</h4>
            <div class="space-y-1">
              <div class="grid grid-cols-[40px_1fr_60px_70px] gap-2 text-xs font-medium text-gray-500 pb-2 border-b border-gray-100">
                <span>排名</span>
                <span>用户</span>
                <span>消息数</span>
                <span>活跃度</span>
              </div>
              {#each groupData.topUsers as user, index}
                <div class="grid grid-cols-[40px_1fr_60px_70px] gap-2 items-center py-2 hover:bg-gray-50 rounded-lg px-2 -mx-2">
                  <span class="text-sm font-medium text-gray-700">{index + 1}</span>
                  <div class="flex items-center gap-2">
                    <div class="w-6 h-6 rounded-full {index === 0 ? 'bg-green-500' : index === 1 ? 'bg-blue-500' : index === 2 ? 'bg-orange-500' : index === 3 ? 'bg-red-500' : 'bg-purple-500'}"></div>
                    <span class="text-sm font-medium text-gray-900 truncate">{user.username}</span>
                  </div>
                  <span class="text-sm font-semibold text-gray-700">{user.messages}</span>
                  <span class="text-sm font-semibold text-green-600">{user.activity}%</span>
                </div>
              {/each}
            </div>
          </div>
          
          <!-- 24小时活跃度分布 -->
          <div class="bg-white border border-gray-200 rounded-xl p-6">
            <h4 class="text-lg font-semibold text-gray-900 mb-4">24小时活跃度分布</h4>
            <div class="mb-4">
              <!-- 简化的图表显示 -->
              <div class="flex items-end justify-between h-32 mb-2">
                {#each groupData.hourlyActivity as activity, hour}
                  <div 
                    class="bg-gradient-to-t {activity > 0.7 ? 'from-red-500 to-red-300' : activity > 0.4 ? 'from-orange-500 to-orange-300' : 'from-green-500 to-green-300'} w-2 rounded-t" 
                    style="height: {activity * 100}%"
                    title="{hour}:00 - 活跃度: {Math.round(activity * 100)}%"
                  ></div>
                {/each}
              </div>
              <div class="flex justify-between text-xs text-gray-500">
                <span>0</span>
                <span>6</span>
                <span>12</span>
                <span>18</span>
              </div>
            </div>
            <!-- 图例 -->
            <div class="flex gap-4 text-xs">
              <div class="flex items-center gap-1">
                <div class="w-3 h-3 bg-red-500 rounded"></div>
                <span class="text-gray-600">高峰期</span>
              </div>
              <div class="flex items-center gap-1">
                <div class="w-3 h-3 bg-orange-500 rounded"></div>
                <span class="text-gray-600">活跃期</span>
              </div>
              <div class="flex items-center gap-1">
                <div class="w-3 h-3 bg-green-500 rounded"></div>
                <span class="text-gray-600">低谷期</span>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 底部信息 -->
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 pt-6 border-t border-gray-200">
          <div class="text-sm text-gray-500">
            <span>数据更新时间: {groupData.lastUpdate}</span>
          </div>
          <div class="flex gap-3">
            <button 
              class="px-6 py-2 bg-indigo-600 text-white font-medium rounded-lg hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 transition-colors" 
              on:click={closeModal}
            >
              加入群组
            </button>
            <button 
              class="px-6 py-2 bg-gray-100 text-gray-700 font-medium rounded-lg hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2 transition-colors" 
              on:click={closeModal}
            >
              关闭
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
{/if}



<!--
  关键算法说明：
  - 组件采用事件派发机制与父组件通信
  - 支持ESC键和遮罩层点击关闭弹窗
  - 响应式设计适配移动端和桌面端
  待优化事项：
  - 可添加更多图表类型支持
  - 支持数据导出功能
  - 添加更多交互动画效果
  兼容性说明：
  - 兼容现代浏览器，使用CSS Grid和Flexbox布局
  - 支持触摸设备的交互体验
-->