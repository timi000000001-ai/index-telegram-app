<!--
 * Telegram搜索平台 - 每日动态页面
 * @description 展示群组、频道的每日动态信息，包括人数排名和分类统计
 * @author 前端工程师
 * @date 2024-01-23
 * @version 1.0.0
 * © Telegram搜索平台
-->

<script>
  /**
   * 每日动态页面逻辑
   * @description 处理动态数据展示和tab切换
   */
  import { onMount } from 'svelte';
  import { apiFetch } from '$lib/api.js';
  
  /**
   * 当前选中的tab
   * @type {string} 'ranking' | 'category'
   */
  let activeTab = 'ranking';
  
  /**
   * 动态数据
   * @type {Object[]}
   */
  let statsData = [];
  
  /**
   * 加载状态
   * @type {boolean}
   */
  let loading = true;
  
  /**
   * 错误信息
   * @type {string}
   */
  let error = '';
  
  /**
   * 获取每日动态数据
   * @description 从API获取群组、频道的动态统计信息
   * @returns {Promise<void>}
   */
  async function fetchDailyStats() {
    try {
      loading = true;
      error = '';
      
      // 模拟API调用 - 实际项目中替换为真实API
      const response = await new Promise(resolve => {
        setTimeout(() => {
          resolve({
            success: true,
            data: generateMockData()
          });
        }, 1000);
      });
      
      if (response.success) {
        statsData = response.data;
      } else {
        error = '获取数据失败';
      }
    } catch (err) {
      error = '网络错误，请稍后重试';
      console.error('获取每日动态数据失败:', err);
    } finally {
      loading = false;
    }
  }
  
  /**
   * 生成模拟数据
   * @description 生成用于展示的模拟动态数据
   * @returns {Object[]} 模拟的动态数据
   */
  function generateMockData() {
    const types = ['群组', '频道'];
    const categories = ['科技', '娱乐', '新闻', '教育', '游戏', '生活'];
    
    return Array.from({ length: 20 }, (_, i) => ({
      id: i + 1,
      name: `${types[i % 2]}${i + 1}`,
      type: types[i % 2],
      category: categories[i % categories.length],
      totalMembers: Math.floor(Math.random() * 50000) + 1000,
      newMembers: Math.floor(Math.random() * 500) + 10,
      messages24h: Math.floor(Math.random() * 1000) + 50,
      activeUsers: Math.floor(Math.random() * 200) + 20,
      dailyScore: (Math.random() * 10).toFixed(1),
      avatar: `https://api.dicebear.com/7.x/initials/svg?seed=${encodeURIComponent('群组' + (i + 1))}`
    })).sort((a, b) => Number(b.totalMembers) - Number(a.totalMembers));
  }
  
  /**
   * 切换tab
   * @param {string} tab - 要切换到的tab
   */
  function switchTab(tab) {
    activeTab = tab;
  }
  
  /**
   * 按分类分组数据
   * @param {Object[]} data - 原始数据
   * @returns {Object} 按分类分组的数据
   */
  function groupByCategory(data) {
    return data.reduce((acc, item) => {
      if (!acc[item.category]) {
        acc[item.category] = [];
      }
      acc[item.category].push(item);
      return acc;
    }, {});
  }
  
  /**
   * 页面初始化
   */
  onMount(() => {
    fetchDailyStats();
  });
  
  // 响应式计算
  $: groupedData = groupByCategory(statsData);
</script>

<!--
 * 页面标题和描述
-->
<svelte:head>
  <title>每日动态 - Telegram搜索平台</title>
  <meta name="description" content="查看Telegram群组和频道的每日动态统计，包括人数排名、活跃度等信息" />
</svelte:head>

<!--
 * 页面主体内容
-->
<div class="max-w-7xl mx-auto">
  <!-- 页面标题 -->
  <div class="mb-8">
    <h1 class="text-3xl font-bold text-slate-800 mb-2">每日动态</h1>
    <p class="text-slate-600">实时统计群组、频道的动态信息和活跃度排名</p>
  </div>
  
  <!-- Tab导航 -->
  <div class="bg-white rounded-lg shadow-sm border border-slate-200 mb-6">
    <div class="flex border-b border-slate-200">
      <button 
        class="tab-button {activeTab === 'ranking' ? 'tab-active' : ''}"
        on:click={() => switchTab('ranking')}>
        <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path>
        </svg>
        人数排名
      </button>
      <button 
        class="tab-button {activeTab === 'category' ? 'tab-active' : ''}"
        on:click={() => switchTab('category')}>
        <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14-4H5m14 8H5m14 4H5"></path>
        </svg>
        分类统计
      </button>
    </div>
  </div>
  
  <!-- 内容区域 -->
  <div class="bg-white rounded-lg shadow-sm border border-slate-200">
    {#if loading}
      <!-- 加载状态 -->
      <div class="flex items-center justify-center py-12">
        <div class="flex items-center gap-3 text-slate-500">
          <svg class="w-6 h-6 animate-spin" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="m4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <span>加载中...</span>
        </div>
      </div>
    {:else if error}
      <!-- 错误状态 -->
      <div class="flex flex-col items-center justify-center py-12">
        <svg class="w-12 h-12 text-red-500 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
        </svg>
        <p class="text-slate-600 mb-4">{error}</p>
        <button 
          class="bg-blue-500 text-white px-4 py-2 rounded-lg hover:bg-blue-600 transition-colors"
          on:click={fetchDailyStats}>
          重新加载
        </button>
      </div>
    {:else}
      <!-- 数据展示 -->
      {#if activeTab === 'ranking'}
        <!-- 人数排名tab -->
        <div class="p-6">
          <div class="grid gap-4">
            {#each statsData as item, index}
              <div class="flex flex-col md:flex-row items-start md:items-center gap-3 md:gap-4 p-3 md:p-4 bg-slate-50 rounded-lg hover:bg-slate-100 transition-colors">
                
                
                <!-- 头像 -->
                <img src={item.avatar} alt={item.name} class="w-12 h-12 rounded-full" />
                
                <!-- 基本信息 -->
                <div class="flex-1">
                  <div class="flex flex-wrap items-center gap-2 mb-1 text-sm md:text-base">
                    <h3 class="font-semibold text-slate-800">{item.name}</h3>
                    <span class="px-2 py-1 bg-blue-100 text-blue-700 text-xs rounded-full">{item.type}</span>
                    <span class="px-2 py-1 bg-green-100 text-green-700 text-xs rounded-full">{item.category}</span>
                  </div>
                  <div class="flex flex-col md:flex-row flex-wrap items-start md:items-center gap-2 md:gap-4 text-xs md:text-sm text-slate-600">
                    <span>总成员: <strong class="text-slate-800">{item.totalMembers.toLocaleString()}</strong></span>
                    <span>新增: <strong class="text-green-600">+{item.newMembers}</strong></span>
                    <span>24h消息: <strong class="text-blue-600">{item.messages24h}</strong></span>
                    <span>活跃用户: <strong class="text-purple-600">{item.activeUsers}</strong></span>
                  </div>
                </div>
                
                <!-- 评分 -->
                <div class="flex-shrink-0 flex items-center gap-2 text-right">
                  <div class="text-lg font-bold text-amber-600">{item.dailyScore}</div>
                  <div class="text-xs text-slate-500">当日评分</div>
                </div>
              </div>
            {/each}
          </div>
        </div>
      {:else}
        <!-- 分类统计tab -->
        <div class="p-6">
          <div class="grid gap-6">
            {#each Object.entries(groupedData) as [category, items]}
              <div class="border border-slate-200 rounded-lg overflow-hidden">
                <div class="bg-slate-100 px-4 py-3 border-b border-slate-200">
                  <h3 class="font-semibold text-slate-800 flex items-center gap-2">
                    <span class="w-3 h-3 bg-gradient-to-r from-blue-500 to-purple-600 rounded-full"></span>
                    {category} ({items.length})
                  </h3>
                </div>
                <div class="divide-y divide-slate-200">
                  {#each items as item}
                    <div class="flex flex-col md:flex-row items-start md:items-center gap-3 md:gap-4 p-3 md:p-4 hover:bg-slate-50 transition-colors">
                      <img src={item.avatar} alt={item.name} class="w-10 h-10 rounded-full" />
                      <div class="flex-1">
                        <div class="flex flex-wrap items-center gap-2 mb-1 text-sm md:text-base">
                          <h4 class="font-medium text-slate-800">{item.name}</h4>
                          <span class="px-2 py-1 bg-blue-100 text-blue-700 text-xs rounded-full">{item.type}</span>
                        </div>
                        <div class="flex flex-col md:flex-row flex-wrap items-start md:items-center gap-2 md:gap-3 text-xs md:text-sm text-slate-600">
                          <span>{item.totalMembers.toLocaleString()} 成员</span>
                          <span class="text-green-600">+{item.newMembers} 新增</span>
                          <span class="text-amber-600">评分 {item.dailyScore}</span>
                        </div>
                      </div>
                    </div>
                  {/each}
                </div>
              </div>
            {/each}
          </div>
        </div>
      {/if}
    {/if}
  </div>
</div>

<!--
 * 样式定义
-->
<style>
  /* Tab按钮样式 */
  .tab-button {
    @apply flex items-center px-6 py-4 text-slate-600 font-medium transition-all duration-200;
    @apply hover:text-blue-600 hover:bg-blue-50;
  }
  
  .tab-active {
    @apply text-blue-600 bg-blue-50 border-b-2 border-blue-600;
  }
  
  /* 响应式设计 */
  @media (max-width: 768px) {
    .tab-button {
      @apply px-4 py-3 text-sm;
    }
  }
</style>

<!--
 * 文件说明
 * @description 每日动态页面，展示群组和频道的实时统计信息
 * @features 人数排名、分类统计、实时数据、响应式设计
 * @compatibility 支持现代浏览器，兼容移动端
-->