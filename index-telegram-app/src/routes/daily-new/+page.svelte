<!--
 * Telegram搜索平台 - 每日新增页面
 * @description 展示新收录的群组、频道、机器人信息
 * @author 前端工程师
 * @date 2024-01-23
 * @version 1.0.0
 * © Telegram搜索平台
-->

<script>
  /**
   * 每日新增页面逻辑
   * @description 处理新增数据展示和tab切换
   */
  import { onMount } from 'svelte';
import { apiFetch } from '$lib/api.js';
import GroupModal from '$lib/components/GroupModal.svelte';
  import { writable } from 'svelte/store';

  const collapsedCategories = writable({});
  
  /**
   * 当前选中的tab
   * @type {string} 'members' | 'category'
   */
  let activeTab = 'members';
  
  /**
   * 新增数据
   * @type {Object[]}
   */
  let newData = [];
  
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
let showGroupModal = false;
let selectedGroupData = null;
  
  /**
   * 获取每日新增数据
   * @description 从API获取新收录的群组、频道、机器人信息
   * @returns {Promise<void>}
   */
  async function fetchDailyNew() {
    try {
      loading = true;
      error = '';
      
      // 模拟API调用 - 实际项目中替换为真实API
      const response = await new Promise(resolve => {
        setTimeout(() => {
          resolve({
            success: true,
            data: generateMockNewData()
          });
        }, 1000);
      });
      
      if (response.success) {
        newData = response.data;
      } else {
        error = '获取数据失败';
      }
    } catch (err) {
      error = '网络错误，请稍后重试';
      console.error('获取每日新增数据失败:', err);
    } finally {
      loading = false;
    }
  }
  
  /**
   * 生成模拟新增数据
   * @description 生成用于展示的模拟新增数据
   * @returns {Object[]} 模拟的新增数据
   */
  function generateMockNewData() {
    const types = ['群组', '频道', '机器人'];
    const categories = ['科技', '娱乐', '新闻', '教育', '游戏', '生活', '工具', '商业'];
    
    return Array.from({ length: 30 }, (_, i) => {
      const type = types[i % 3];
      const addedTime = new Date(Date.now() - Math.random() * 24 * 60 * 60 * 1000);
      
      return {
        id: i + 1,
        name: `${type}${i + 1}`,
        type: type,
        category: categories[i % categories.length],
        description: `这是一个关于${categories[i % categories.length]}的${type}，提供优质内容和服务。`,
        members: type === '机器人' ? 0 : Math.floor(Math.random() * 10000) + 100,
        addedTime: addedTime,
        addedTimeStr: formatTime(addedTime),
        isVerified: Math.random() > 0.7,
        isNew: Math.random() > 0.5,
        avatar: `https://api.dicebear.com/7.x/initials/svg?seed=${encodeURIComponent(type + (i + 1))}`,
        tags: generateTags(categories[i % categories.length], type)
      };
    }).sort((a, b) => b.addedTime - a.addedTime);
  }
  
  /**
   * 生成标签
   * @param {string} category - 分类
   * @param {string} type - 类型
   * @returns {string[]} 标签数组
   */
  function generateTags(category, type) {
    const tagMap = {
      '科技': ['AI', '编程', '开发', '技术'],
      '娱乐': ['音乐', '电影', '游戏', '综艺'],
      '新闻': ['时事', '财经', '国际', '本地'],
      '教育': ['学习', '考试', '培训', '知识'],
      '游戏': ['手游', '端游', '攻略', '竞技'],
      '生活': ['美食', '旅行', '健康', '购物'],
      '工具': ['效率', '实用', '办公', '助手'],
      '商业': ['创业', '投资', '营销', '管理']
    };
    
    const baseTags = tagMap[category] || ['其他'];
    const selectedTags = baseTags.slice(0, Math.floor(Math.random() * 3) + 1);
    
    if (type === '机器人') {
      selectedTags.push('自动化');
    }
    
    return selectedTags;
  }
  
  /**
   * 格式化时间
   * @param {Date} date - 日期对象
   * @returns {string} 格式化的时间字符串
   */
  function formatTime(date) {
    const now = new Date();
    const diff = now - date;
    const hours = Math.floor(diff / (1000 * 60 * 60));
    const minutes = Math.floor(diff / (1000 * 60));
    
    if (hours < 1) {
      return `${minutes}分钟前`;
    } else if (hours < 24) {
      return `${hours}小时前`;
    } else {
      return date.toLocaleDateString('zh-CN', {
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      });
    }
  }
  
  /**
   * 切换tab
   * @param {string} tab - 要切换到的tab
   */
  function switchTab(tab) {
    activeTab = tab;
  }

function openGroupModal(item) {
  selectedGroupData = {
    name: item.name,
    username: item.username || '',
    description: item.description,
    createdAt: item.addedTimeStr,
    stats: {
      totalMembers: item.members,
      onlineMembers: Math.floor(item.members * 0.1),
      messages24h: Math.floor(Math.random() * 1000) + 500,
      avgDaily: Math.floor(Math.random() * 100) + 50,
      activityRate: Math.floor(Math.random() * 100),
      groupScore: (Math.random() * 10).toFixed(1)
    },
    topUsers: [
      { username: 'user1', messages: 100, activity: 80 },
      { username: 'user2', messages: 90, activity: 75 },
      { username: 'user3', messages: 80, activity: 70 },
      { username: 'user4', messages: 70, activity: 65 },
      { username: 'user5', messages: 60, activity: 60 }
    ],
    hourlyActivity: Array(24).fill(0).map(() => Math.random()),
    lastUpdate: new Date().toLocaleString()
  };
  showGroupModal = true;
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
   * 按人数排序数据
   * @param {Object[]} data - 原始数据
   * @returns {Object[]} 按人数排序的数据
   */
  function sortByMembers(data) {
    return [...data].sort((a, b) => Number(b.members) - Number(a.members));
  }
  
  /**
   * 页面初始化
   */
  onMount(() => {
    fetchDailyNew();
  });
  
  // 响应式计算
  $: groupedData = groupByCategory(newData);
  $: sortedData = sortByMembers(newData);

  function toggleCategory(category) {
    collapsedCategories.update(current => ({
      ...current,
      [category]: !current[category]
    }));
  }
</script>

<!--
 * 页面标题和描述
-->
<svelte:head>
  <title>每日新增 - Telegram搜索平台</title>
  <meta name="description" content="查看最新收录的Telegram群组、频道和机器人" />
</svelte:head>

<!--
 * 页面主体内容
-->
<div class="max-w-7xl mx-auto">
  <!-- 页面标题 -->
  <div class="mb-8">
    <h1 class="text-3xl font-bold text-slate-800 mb-2">每日新增</h1>
    <p class="text-slate-600">发现最新收录的优质群组、频道和机器人</p>
  </div>
  
  <!-- Tab导航 -->
  <div class="bg-white rounded-lg shadow-sm border border-slate-200 mb-6">
    <div class="flex border-b border-slate-200">
      <button 
        class="tab-button {activeTab === 'members' ? 'tab-active' : ''}"
        on:click={() => switchTab('members')}>
        <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"></path>
        </svg>
        按人数排序
      </button>
      <button 
        class="tab-button {activeTab === 'category' ? 'tab-active' : ''}"
        on:click={() => switchTab('category')}>
        <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z"></path>
        </svg>
        按分类浏览
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
          on:click={fetchDailyNew}>
          重新加载
        </button>
      </div>
    {:else}
      <!-- 数据展示 -->
      {#if activeTab === 'members'}
        <!-- 按人数排序tab -->
        <div class="p-6">
          <div class="grid gap-4">
            {#each sortedData as item}
              <div class="flex flex-col md:flex-row items-start md:items-center gap-4 p-4 bg-slate-50 rounded-lg hover:bg-slate-100 transition-colors">
                <!-- 头像 -->
                <img src={item.avatar} alt={item.name} class="w-12 h-12 rounded-full flex-shrink-0" />
                
                <!-- 主要信息 -->
                <div class="flex-1 min-w-0">
                  <div class="flex items-center gap-2 mb-2">
                    <h3 class="font-semibold text-slate-800 truncate">{item.name}</h3>
                    <span class="px-2 py-1 bg-blue-100 text-blue-700 text-xs rounded-full flex-shrink-0">{item.type}</span>
                    {#if item.isVerified}
                      <svg class="w-4 h-4 text-blue-500 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                        <path fill-rule="evenodd" d="M6.267 3.455a3.066 3.066 0 001.745-.723 3.066 3.066 0 013.976 0 3.066 3.066 0 001.745.723 3.066 3.066 0 012.812 2.812c.051.643.304 1.254.723 1.745a3.066 3.066 0 010 3.976 3.066 3.066 0 00-.723 1.745 3.066 3.066 0 01-2.812 2.812 3.066 3.066 0 00-1.745.723 3.066 3.066 0 01-3.976 0 3.066 3.066 0 00-1.745-.723 3.066 3.066 0 01-2.812-2.812 3.066 3.066 0 00-.723-1.745 3.066 3.066 0 010-3.976 3.066 3.066 0 00.723-1.745 3.066 3.066 0 012.812-2.812zm7.44 5.252a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"></path>
                      </svg>
                    {/if}
                    {#if item.isNew}
                      <span class="px-2 py-1 bg-red-100 text-red-700 text-xs rounded-full flex-shrink-0">NEW</span>
                    {/if}
                  </div>
                  
                  <p class="text-slate-600 text-sm mb-2 line-clamp-2">{item.description}</p>
                  
                  <div class="flex flex-wrap items-center gap-4 text-sm text-slate-500 mb-2">
                    <span class="flex items-center gap-1">
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"></path>
                      </svg>
                      {item.type === '机器人' ? '机器人' : `${item.members.toLocaleString()} 成员`}
                    </span>
                    <span class="flex items-center gap-1">
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                      </svg>
                      {item.addedTimeStr}
                    </span>
                    <span class="flex items-center gap-1">
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z"></path>
                      </svg>
                      {item.category}
                    </span>
                  </div>
                  
                  <!-- 标签 -->
                  <div class="flex flex-wrap gap-1">
                    {#each item.tags as tag}
                      <span class="px-2 py-1 bg-gray-100 text-gray-600 text-xs rounded">{tag}</span>
                    {/each}
                  </div>
                </div>
                
                <!-- 操作按钮 -->
                <div class="flex flex-col gap-2 flex-shrink-0">
                  <button class="px-3 py-1 bg-blue-500 text-white text-sm rounded hover:bg-blue-600 transition-colors" on:click={() => openGroupModal(item)}>
                    查看详情
                  </button>
                  <button class="px-3 py-1 bg-gray-100 text-gray-600 text-sm rounded hover:bg-gray-200 transition-colors">
                    收藏
                  </button>
                </div>
              </div>
            {/each}
          </div>
        </div>
      {:else}
        <!-- 按分类浏览tab -->
        <div class="p-6">
          <div class="grid gap-6">
            {#each Object.entries(groupedData) as [category, items]}
              <div class="border border-slate-200 rounded-lg overflow-hidden">
                <div class="bg-slate-100 px-4 py-3 border-b border-slate-200 cursor-pointer flex justify-between items-center" on:click={() => toggleCategory(category)}>
                  <h3 class="font-semibold text-slate-800 flex items-center gap-2">
                    <span class="w-3 h-3 bg-gradient-to-r from-blue-500 to-purple-600 rounded-full"></span>
                    {category} ({items.length})
                  </h3>
                  <svg
                    class="w-5 h-5 text-slate-500 transform transition-transform duration-200"
                    class:rotate-180={!$collapsedCategories[category]}
                    fill="none" stroke="currentColor" viewBox="0 0 24 24"
                  >
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
                  </svg>
                </div>
                {#if !$collapsedCategories[category]}
                  <div class="divide-y divide-slate-200">
                    {#each items as item}
                    <div class="flex flex-col md:flex-row items-start md:items-center gap-4 p-4 hover:bg-slate-50 transition-colors">
                      <img src={item.avatar} alt={item.name} class="w-10 h-10 rounded-full" />
                      <div class="flex-1 min-w-0">
                        <div class="flex items-center gap-2 mb-1">
                          <h4 class="font-medium text-slate-800 truncate">{item.name}</h4>
                          <span class="px-2 py-1 bg-blue-100 text-blue-700 text-xs rounded-full">{item.type}</span>
                          {#if item.isVerified}
                            <svg class="w-4 h-4 text-blue-500" fill="currentColor" viewBox="0 0 20 20">
                              <path fill-rule="evenodd" d="M6.267 3.455a3.066 3.066 0 001.745-.723 3.066 3.066 0 013.976 0 3.066 3.066 0 001.745.723 3.066 3.066 0 012.812 2.812c.051.643.304 1.254.723 1.745a3.066 3.066 0 010 3.976 3.066 3.066 0 00-.723 1.745 3.066 3.066 0 01-2.812 2.812 3.066 3.066 0 00-1.745.723 3.066 3.066 0 01-3.976 0 3.066 3.066 0 00-1.745-.723 3.066 3.066 0 01-2.812-2.812 3.066 3.066 0 00-.723-1.745 3.066 3.066 0 010-3.976 3.066 3.066 0 00.723-1.745 3.066 3.066 0 012.812-2.812zm7.44 5.252a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"></path>
                            </svg>
                          {/if}
                          {#if item.isNew}
                            <span class="px-2 py-1 bg-red-100 text-red-700 text-xs rounded-full">NEW</span>
                          {/if}
                        </div>
                        <div class="flex flex-wrap items-center gap-3 text-sm text-slate-600">
                          <span>{item.type === '机器人' ? '机器人' : `${item.members.toLocaleString()} 成员`}</span>
                          <span>{item.addedTimeStr}</span>
                        </div>
                      </div>
                      <div class="flex gap-2">
                        <button class="px-3 py-1 bg-blue-500 text-white text-sm rounded hover:bg-blue-600 transition-colors" on:click={() => openGroupModal(item)}>
                          查看
                        </button>
                      </div>
                    </div>
                  {/each}
                </div>
                {/if}
              </div>
            {/each}
          </div>
        </div>
      {/if}
    {/if}
  </div>
</div>

<GroupModal 
  show={showGroupModal} 
  groupData={selectedGroupData} 
  on:close={() => showGroupModal = false} 
/>

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
  
  /* 文本截断 */
  .line-clamp-2 {
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
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
 * @description 每日新增页面，展示最新收录的群组、频道和机器人
 * @features 人数排序、分类浏览、实时数据、响应式设计
 * @compatibility 支持现代浏览器，兼容移动端
-->