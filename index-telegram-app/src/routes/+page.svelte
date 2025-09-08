<script>
/**
 * @fileoverview Telegram 内容搜索主页面 - 优化版本
 * @description 提供现代化的搜索界面，包含智能建议、高级筛选和美观的结果展示
 * @author 前端工程师
 * @date 2024-01-20
 * @version 2.0.0
 * © Telegram Search Platform
 */

import { onMount, onDestroy } from 'svelte';
import { browser } from '$app/environment';
import { Meilisearch } from 'meilisearch';
import GroupModal from '$lib/components/GroupModal.svelte';

import CollectionModal from '$lib/components/CollectionModal.svelte';
import DisclaimerModal from '$lib/components/DisclaimerModal.svelte';

// ===================== 状态管理 =====================
/** @type {string} 搜索关键字 */
let query = '';
/** @type {any[]} 搜索结果列表 */
let results = [];
/** @type {string[]} 搜索建议列表 */
let suggestions = [];
/** @type {any[]} 热门搜索关键字 */
let trending = [];
/** @type {string[]} 搜索历史记录 */
let histories = [];


// 搜索状态
/** @type {boolean} 是否正在搜索 */
let isSearching = false;
/** @type {boolean} 是否正在获取建议 */
let isFetchingSuggestions = false;
/** @type {boolean} 是否正在获取热门关键字 */
let isFetchingTrending = false;
/** @type {number} 搜索耗时（毫秒） */
let elapsedMs = 0;

// UI 状态
/** @type {boolean} 是否显示搜索建议 */
let showSuggestions = false;
/** @type {boolean} 是否显示高级选项 */
let showAdvanced = false;
/** @type {boolean} 移动端是否显示筛选器 */
let showMobileFilters = false;
/** @type {boolean} 是否显示群组详情弹窗 */
let showGroupModal = false;
/** @type {boolean} 是否显示收录弹窗 */
let showCollectionModal = false;

// 文本动画控制

/** @type {any} 当前选中的群组/频道详情数据 */
let selectedGroupData = null;

// 分页与排序
/** @type {number} 当前页码 */
let page = 1;
/** @type {number} 每页条数 */
let size = 10;
/** @type {number} 总结果数 */
let total = 0;
/** @type {string} 排序方式 */
let sort = 'relevance';

// 高级筛选器
/** @type {any} */
let filters = {
  types: {
    group: true,
    channel: true,
    bot: true,
    message: true
  },
  timePreset: 'any',
  timeOrder: 'desc',
  customStart: '',
  customEnd: '',
  sourceGroup: '',
  sourceChannel: '',
  sourceUser: '',
  language: '',
  minLength: '',
  minInteractions: '',
  minRelevance: ''
};

// ===================== 工具函数 =====================
/**
 * 防抖函数
 * @param {Function} fn - 要防抖的函数
 * @param {number} wait - 等待时间（毫秒）
 * @returns {Function} 防抖后的函数
 */
function debounce(fn, wait) {
   /** @type {number | undefined} */
   let timeout;
   return function executedFunction(/** @type {...any} */ ...args) {
     const later = () => {
       clearTimeout(timeout);
       fn(...args);
     };
     clearTimeout(timeout);
     timeout = setTimeout(later, wait);
   };
  }

/**
 * 高亮搜索关键字
 * @param {string} text - 原始文本
 * @param {string} query - 搜索关键字
 * @returns {string} 高亮后的HTML
 */
function highlight(text, query) {
  if (!query.trim()) return text;
  const keywords = query.trim().split(/\s+/);
  let result = text;
  keywords.forEach(keyword => {
    const regex = new RegExp(`(${keyword})`, 'gi');
    result = result.replace(regex, '<mark>$1</mark>');
  });
  
  return result;
}

/**
 * 计算总页数
 * @returns {number} 总页数
 */
function totalPages() {
  return Math.ceil(total / size);
}

const client = new Meilisearch({
  host: 'http://127.0.0.1:7700',
  apiKey: 'timigogogo',
});

const index = client.index('telegram_index');

// ===================== API 调用 =====================


/**
 * 执行搜索
 */
async function search() {
  if (!query.trim()) {
    results = [];
    total = 0;
    return;
  }

  isSearching = true;
  const startTime = Date.now();

  try {
    const searchResults = await index.search(query.trim(), {
      page: page,
      hitsPerPage: size,
      filter: [`TYPE IN [${Object.entries(filters.types).filter(([_,v]) => v).map(([k,_])=>`'${k}'`).join(', ')}]`]
    });
    results = searchResults.hits.map(hit => {
      /** @type {{[key: string]: any}} */
      const newHit = {};
      for (const key in hit) {
        newHit[key.toLowerCase()] = /** @type {any} */ (hit)[key];
      }
      if (newHit.link && typeof newHit.link === 'string') {
        newHit.link = newHit.link.replace(/`/g, '').trim();
      }
      return newHit;
    });
    total = searchResults.totalHits ?? 0;
    elapsedMs = searchResults.processingTimeMs;

  } catch (error) {
    console.error('搜索失败:', error);
    results = [];
    total = 0;
  } finally {
    isSearching = false;
  }
}

/**
 * 获取搜索建议
 * @param {string} currentQuery - 当前查询
 * @returns {Promise<void>}
 */
async function fetchSuggestions(currentQuery) {
  if (!currentQuery || currentQuery.length < 2) {
    suggestions = [];
    showSuggestions = true; // 即使没有建议，也要确保建议列表是“活动的”，以便显示历史记录或热门
    return;
  }

  isFetchingSuggestions = true;
  try {
    const response = await fetch(`/api/search/autocomplete?q=${encodeURIComponent(currentQuery)}`);
    if (response.ok) {
      const data = await response.json();
      suggestions = data;
    } else {
      console.error('Failed to fetch suggestions');
      suggestions = [];
    }
  } catch (error) {
    console.error('Error fetching suggestions:', error);
    suggestions = [];
  } finally {
    isFetchingSuggestions = false;
    showSuggestions = true;
  }
}

/**
 * 获取热门搜索
 */
async function fetchTrending() {
  isFetchingTrending = true;
  try {
    // This is a mock API, replace with your actual implementation
    const mockTrending = [
      { keyword: 'SvelteKit', rank: 1, count: 120, category: 'technology', trend: 'up' },
      { keyword: 'Meilisearch', rank: 2, count: 105, category: 'technology', trend: 'hot' },
      { keyword: 'TailwindCSS', rank: 3, count: 98, category: 'technology', trend: 'stable' },
    ];
    await new Promise(resolve => setTimeout(resolve, 500));
    trending = mockTrending;

  } catch (error) {
    console.error('获取热门搜索失败:', error);
    trending = [];
  } finally {
    isFetchingTrending = false;
  }
}

// ===================== 辅助函数 =====================
/**
 * 获取分类中文名称
 * @param {string} category - 分类英文名
 * @returns {string} 分类中文名
 */
function getCategoryName(category) {
  /** @type {Record<string, string>} */
  const categoryMap = {
    technology: '技术',
    blockchain: '区块链',
    ai: '人工智能',
    mobile: '移动开发',
    devops: '运维开发'
  };
  return categoryMap[category] || '其他';
}

/**
 * 获取趋势中文名称
 * @param {string} trend - 趋势英文名
 * @returns {string} 趋势中文名
 */
function getTrendName(trend) {
  /** @type {Record<string, string>} */
  const trendMap = {
    up: '上升',
    down: '下降',
    hot: '热门',
    stable: '稳定'
  };
  return trendMap[trend] || '未知';
}

/**
 * 获取趋势图标
 * @param {string} trend - 趋势类型
 * @returns {string} 趋势图标
 */
function getTrendIcon(trend) {
  /** @type {Record<string, string>} */
  const iconMap = {
    up: '📈',
    down: '📉',
    hot: '🔥',
    stable: '➡️'
  };
  return iconMap[trend] || '❓';
}

// ===================== 事件处理 =====================

const debouncedFetchSuggestions = debounce(fetchSuggestions, 300);

/**
 * 处理搜索输入变化
 * @param {Event} e - 输入事件
 * @returns {void}
 */
function onQueryInput(e) {
  const target = /** @type {HTMLInputElement} */ (e.target);
  query = target.value;
  debouncedFetchSuggestions(query);
}

/**
 * 选择搜索建议
 * @param {string} suggestion - 选择的建议
 * @returns {void}
 */
function selectSuggestion(suggestion) {
  query = suggestion;
  suggestions = [];
  showSuggestions = false;
  page = 1;
  search();
  saveToHistory(query);
}

/**
   * 切换高级选项显示
   * @returns {void}
   */
  function toggleAdvanced() {
    showAdvanced = !showAdvanced;
  }

  /**
   * 切换移动端筛选器显示
   * @returns {void}
   */
  function toggleMobileFilters() {
    showMobileFilters = !showMobileFilters;
  }

/**
 * 跳转到指定页面
 * @param {number} targetPage - 目标页码
 * @returns {void}
 */
function goToPage(targetPage) {
  if (targetPage < 1 || targetPage > totalPages()) return;
  page = targetPage;
  search();
}

/**
 * 重置筛选器
 * @returns {void}
 */
function resetFilters() {
  filters = {
    types: {
      group: true,
      channel: true,
      bot: true,
      message: true
    },
    timePreset: 'any',
    timeOrder: 'desc',
    customStart: '',
    customEnd: '',
    sourceGroup: '',
    sourceChannel: '',
    sourceUser: '',
    language: '',
    minLength: '',
    minInteractions: '',
    minRelevance: ''
  };
}

/**
 * 保存搜索历史
 * @param {string} searchQuery - 搜索查询
 * @returns {void}
 */
function saveToHistory(searchQuery) {
  if (!browser || !searchQuery) return;
  
  /** @type {string[]} */
  let savedHistories = [];
  try {
    savedHistories = JSON.parse(localStorage.getItem('searchHistories') || '[]');
  } catch (e) {
    savedHistories = [];
  }
  
  // 移除重复项并添加到开头
  savedHistories = savedHistories.filter(h => h !== searchQuery);
  savedHistories.unshift(searchQuery);
  
  // 限制历史记录数量
  savedHistories = savedHistories.slice(0, 10);
  
  localStorage.setItem('searchHistories', JSON.stringify(savedHistories));
  histories = savedHistories;
}

/**
 * 打开群组详情弹窗
 * @param {any} item - 群组或频道数据
 * @returns {void}
 */
function openGroupModal(item) {
  selectedGroupData = generateGroupStatistics(item);
  showGroupModal = true;
}

/**
 * 关闭群组详情弹窗
 * @returns {void}
 */
function closeGroupModal() {
  showGroupModal = false;
  selectedGroupData = null;
}

/**
 * 生成群组统计数据
 * @param {any} item - 群组或频道基础数据
 * @returns {any} 群组统计详情数据
 */
function generateGroupStatistics(item) {
  const isChannel = item.type === 'channel';
  const baseMembers = Math.floor(Math.random() * 2000) + 500;
  const activeUsers = Math.floor(baseMembers * (0.3 + Math.random() * 0.4));
  
  return {
    id: item.id,
    name: item.source,
    type: item.type,
    description: isChannel ? '专业的Telegram机器人开发交流群组，汇聚全球开发者' : '专业的Telegram机器人开发交流频道，汇聚全球开发者',
    username: '@developers_chat',
    createdAt: '2023-05-15',
    
    // 统计数据
    stats: {
      totalMembers: baseMembers,
      onlineMembers: Math.floor(activeUsers * 0.4),
      messages24h: Math.floor(Math.random() * 500) + 100,
      avgDaily: Math.floor(Math.random() * 10) + 3,
      activityRate: Math.floor((activeUsers / baseMembers) * 100),
      groupScore: (Math.random() * 2 + 3).toFixed(1)
    },
    // 活跃用户TOP5
    topUsers: [
      { username: '@developer_alice', messages: 45, activity: 92 },
      { username: '@bot_master', messages: 38, activity: 87 },
      { username: '@api_expert', messages: 32, activity: 81 },
      { username: '@code_ninja', messages: 28, activity: 76 },
      { username: '@tech_guru', messages: 24, activity: 72 }
    ],
    
    // 24小时活跃度分布数据
    hourlyActivity: Array.from({ length: 24 }, (_, i) => {
      const baseActivity = Math.sin((i - 6) * Math.PI / 12) * 0.5 + 0.5;
      const randomFactor = Math.random() * 0.3;
      return Math.max(0.1, Math.min(1, baseActivity + randomFactor));
    }),
    
    // 更新时间
    lastUpdate: new Date().toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit', 
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    })
  };
}

/**
 * 加载搜索历史
 * @returns {void}
 */
function loadHistory() {
  if (!browser) return;
  
  try {
    histories = JSON.parse(localStorage.getItem('searchHistories') || '[]');
  } catch (e) {
    histories = [];
  }
}

/**
 * 清空搜索历史
 * @returns {void}
 */
function clearHistory() {
  if (!browser) return;
  
  localStorage.removeItem('searchHistories');
  histories = [];
}

  let showDisclaimerModal = false;

  function openDisclaimerModal() {
    showDisclaimerModal = true;
  }

  function closeDisclaimerModal() {
    showDisclaimerModal = false;
  }

// ===================== 生命周期 =====================
onMount(() => {
  loadHistory();
  fetchTrending();
  search();
});


</script>

<!-- ===================== 页面结构 ===================== -->

<div class="w-full p-2 min-h-screen bg-gradient-to-br from-slate-50 via-blue-50/30 to-indigo-50/50 relative overflow-hidden">
  
  <!-- 主搜索区域 -->
  <section class="relative z-10 bg-gradient-to-br from-white/80 via-blue-50/60 to-indigo-100/40 backdrop-blur-xl rounded-3xl md:rounded-[2rem] p-6 md:p-8 lg:p-12 shadow-2xl shadow-blue-500/10 border border-white/20 mb-6 overflow-hidden">
    <!-- 搜索区域内的动态背景效果 -->
    <div class="absolute inset-0 pointer-events-none">
      <!-- 所有装饰图标已移除 -->
      
      <!-- 装饰图标已移除 -->
      

      
      <!-- 数据流动效果 - 移除动画 -->
      <div class="absolute top-0 left-0 w-full h-full opacity-40">
        <div class="absolute top-4 left-0 w-0.5 h-16 bg-gradient-to-b from-transparent via-blue-400/50 to-transparent shadow-md shadow-blue-400/30"></div>
        <div class="absolute top-8 right-0 w-0.5 h-12 bg-gradient-to-b from-transparent via-purple-400/40 to-transparent shadow-md shadow-purple-400/20"></div>
        <div class="absolute bottom-8 left-1/4 w-0.5 h-14 bg-gradient-to-b from-transparent via-indigo-400/45 to-transparent shadow-md shadow-indigo-400/25"></div>
        
        <!-- 更多数据流 -->
        <div class="absolute top-12 right-1/3 w-0.5 h-10 bg-gradient-to-b from-transparent via-cyan-400/35 to-transparent shadow-sm shadow-cyan-400/15"></div>
        <div class="absolute bottom-12 left-1/3 w-0.5 h-18 bg-gradient-to-b from-transparent via-rose-400/30 to-transparent shadow-sm shadow-rose-400/10"></div>
      </div>
    </div>
    
    
  
    
    <h1 class="text-3xl md:text-4xl font-bold bg-gradient-to-r from-slate-800 via-blue-700 to-indigo-800 bg-clip-text text-transparent mb-2 relative z-10 drop-shadow-sm">Telegram 内容搜索</h1>
    <p class="text-slate-600/90 text-base md:text-lg mb-6 relative z-10 font-medium">搜索群组、频道和私聊中的消息内容</p>
    
    <!-- 数据库图标已移除 -->
    
    <div class="flex flex-col sm:flex-row gap-4 items-stretch sm:items-center relative {showSuggestions ? 'z-30' : 'z-10'}">
      <div class="relative flex-1 group">
        <!-- 搜索框光晕效果 - 移除动画 -->
        <div class="absolute inset-0 bg-gradient-to-r from-blue-400/30 via-purple-400/30 to-pink-400/30 rounded-full blur-xl opacity-0 group-hover:opacity-100 group-focus-within:opacity-100 transition-all duration-300"></div>
        
        <!-- 主搜索框 - 简化动画 -->
        <input
          class="w-full px-6 py-4 text-base border-2 border-blue-200/50 rounded-full outline-none transition-all duration-300 bg-white/70 backdrop-blur-md text-slate-700 placeholder-slate-400 focus:border-blue-400/80 focus:bg-white/90 focus:shadow-lg hover:border-blue-300/70 hover:bg-white/80"
          placeholder="请输入搜索关键字（支持多个关键字，以空格分隔）"
          bind:value={query}
          on:input={onQueryInput}
          on:focus={() => query.length > 1 && suggestions.length > 0 && (showSuggestions = true)}
          on:blur={() => setTimeout(() => (showSuggestions = false), 200)} 
          on:keydown={(e) => {
            if (e.key === 'Enter') {
              showSuggestions = false;
              search();
              saveToHistory(query);
            }
          }}
        />
        
        <!-- 搜索图标 - 简化动画 -->
        <div class="absolute right-4 top-1/2 -translate-y-1/2 text-slate-400 pointer-events-none transition-colors duration-300 group-focus-within:text-blue-600">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
          </svg>
        </div>
        
        <!-- 搜索建议下拉框 -->
        {#if showSuggestions && (suggestions.length > 0 || isFetchingSuggestions || histories.length > 0)}
          <div class="absolute top-full mt-2 w-full bg-white/90 backdrop-blur-md border border-blue-200/50 rounded-2xl shadow-lg z-50 overflow-hidden">
            <!-- 加载中 -->
            {#if isFetchingSuggestions}
              <div class="p-4 text-center text-slate-500">正在加载建议...</div>
            {/if}

            <!-- 搜索建议 -->
            {#each suggestions as suggestion}
              <div class="px-4 py-2 cursor-pointer hover:bg-blue-100/50" on:mousedown={() => selectSuggestion(suggestion)}>
                {@html suggestion}
              </div>
            {/each}

            <!-- 搜索历史 -->
            {#if !isFetchingSuggestions && suggestions.length === 0 && histories.length > 0}
              <div class="p-2">
                <div class="flex justify-between items-center px-2 py-1">
                  <span class="text-sm font-semibold text-slate-600">搜索历史</span>
                  <button class="text-xs text-blue-500 hover:underline" on:click={clearHistory}>清空</button>
                </div>
                {#each histories as history}
                  <div class="px-2 py-1.5 cursor-pointer hover:bg-blue-100/50 rounded-md" on:mousedown={() => selectSuggestion(history)}>
                    {history}
                  </div>
                {/each}
              </div>
            {/if}
          </div>
        {/if}
      </div>
      
      <button class="relative bg-gradient-to-r from-blue-500 via-indigo-500 to-purple-600 text-white border-0 px-8 py-4 rounded-full cursor-pointer font-semibold transition-all duration-300 shadow-lg hover:shadow-xl disabled:opacity-60 disabled:cursor-not-allowed group" disabled={isSearching} on:click={() => (page=1, search())}>
        <span class="relative z-10">{#if isSearching}搜索中...{:else}搜索{/if}</span>
      </button>
      
      <button class="relative bg-gradient-to-r from-green-500 via-emerald-500 to-teal-600 text-white border-0 px-6 py-4 rounded-full cursor-pointer font-semibold transition-all duration-300 shadow-lg hover:shadow-xl group" on:click={() => showCollectionModal = true}>
        <span class="relative z-10 flex items-center gap-2">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
          </svg>
          收录
        </span>
      </button>
      
      <button class="relative bg-white/60 text-slate-700 border-2 border-blue-200/50 px-6 py-4 rounded-full cursor-pointer transition-all duration-300 backdrop-blur-md hover:bg-white/80 hover:border-blue-300/70 hover:shadow-md focus:outline-none focus:ring-2 focus:ring-blue-400/50 {query.length === 0 ? 'hidden' : 'flex'} md:flex" on:click={toggleAdvanced}>
        <span class="flex items-center gap-2">
          高级选项 
          <span class="transition-transform duration-300 {showAdvanced ? 'rotate-180' : ''}">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
            </svg>
          </span>
        </span>
      </button>
    </div>



    <!-- 热门搜索 -->
    <div class="{query.length === 0 ? 'hidden' : 'flex'} md:flex flex-col gap-6 mt-8 relative z-10">
      <div class="flex items-center gap-3 flex-wrap pb-4">
        <span class="text-black mr-4 font-bold text-sm flex items-center gap-2 flex-shrink-0">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
          </svg>
          历史搜索
        </span>
        {#if histories.length === 0}
          <span class="text-black/60 italic">暂无历史记录</span>
        {:else}
          {#each histories as h}
            <button class="group relative bg-black/10 border border-black/30 px-4 py-2.5 rounded-full cursor-pointer text-black text-sm font-medium transition-all duration-300 backdrop-blur-md hover:bg-black/20 hover:border-black/50 hover:-translate-y-1 hover:shadow-lg hover:shadow-black/20 hover:scale-105 focus:outline-none focus:ring-2 focus:ring-black/50" on:click={() => selectSuggestion(h)}>
              <span class="relative z-10">{h}</span>
              <div class="absolute inset-0 bg-gradient-to-r from-blue-400/20 to-purple-400/20 rounded-full opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
            </button>
          {/each}
          <button class="bg-transparent border-0 text-black/70 cursor-pointer underline transition-all duration-300 hover:text-black hover:scale-105 font-medium text-sm" on:click={clearHistory}>清空历史</button>
        {/if}
      </div>
      <div class="flex items-center gap-3 flex-wrap pb-4">
        <span class="text-black mr-4 font-bold text-sm flex items-center gap-2 flex-shrink-0">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path>
          </svg>
          热门搜索
        </span>
        {#if isFetchingTrending}
          <span class="text-black/60 flex items-center gap-2">
            <div class="w-3 h-3 border border-black/40 border-t-black rounded-full animate-spin"></div>
            加载中...
          </span>
        {:else if trending.length === 0}
          <span class="text-black/60 italic">暂无热门内容</span>
        {:else}
          {#each trending.slice(0, 8) as t}
            <button 
              class="group relative inline-flex items-center gap-2 px-4 py-2.5 rounded-full text-sm font-medium transition-all duration-300 backdrop-blur-md border focus:outline-none focus:ring-2 focus:ring-white/50 {t.trend === 'hot' ? 'bg-gradient-to-r from-red-500/80 to-pink-500/80 text-white border-red-400/50 shadow-lg shadow-red-500/30 hover:from-red-400 hover:to-pink-400 hover:-translate-y-1 hover:shadow-xl hover:shadow-red-500/40 hover:scale-105' : t.trend === 'up' ? 'bg-gradient-to-r from-emerald-500/80 to-green-500/80 text-white border-emerald-400/50 shadow-lg shadow-emerald-500/30 hover:from-emerald-400 hover:to-green-400 hover:-translate-y-1 hover:shadow-xl hover:shadow-emerald-500/40 hover:scale-105' : t.trend === 'down' ? 'bg-gradient-to-r from-slate-500/80 to-gray-500/80 text-white border-slate-400/50 shadow-lg shadow-slate-500/30 hover:from-slate-400 hover:to-gray-400 hover:-translate-y-1 hover:shadow-xl hover:shadow-slate-500/40 hover:scale-105' : 'bg-gradient-to-r from-blue-500/80 to-purple-500/80 text-white border-blue-400/50 shadow-lg shadow-blue-500/30 hover:from-blue-400 hover:to-purple-400 hover:-translate-y-1 hover:shadow-xl hover:shadow-blue-500/40 hover:scale-105'}" 
              title={`排名: #${t.rank} | 搜索次数: ${t.count} | 分类: ${getCategoryName(t.category)} | 趋势: ${getTrendName(t.trend)}`} 
              on:click={() => selectSuggestion(t.keyword)}
            >
              <span class="text-base leading-none animate-bounce">{getTrendIcon(t.trend)}</span>
              <span class="font-bold relative z-10">{t.keyword}</span>
              <span class="text-xs opacity-90 font-medium bg-white/20 px-2 py-1 rounded-full">({t.count})</span>
              <div class="absolute inset-0 bg-gradient-to-r from-white/10 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300 rounded-full"></div>
            </button>
          {/each}
        {/if}
      </div>
    </div>
  </section>

  <!-- 移动端筛选器切换按钮 -->
  <div class="md:hidden mb-4">
    <button class="w-full flex items-center justify-center gap-2 px-4 py-3 bg-blue-600 text-white rounded-2xl font-medium hover:bg-blue-700 transition-colors" on:click={toggleMobileFilters}>
      <span class="text-lg">🔍</span>
      <span>筛选器 {showMobileFilters ? '▲' : '▼'}</span>
    </button>
  </div>

  <!-- 主内容区域 -->
  <div class="flex flex-col md:flex-row gap-3 w-full px-2">
    <!-- 左侧筛选栏 -->
    <aside class="w-full md:w-64 flex-shrink-0 {!showMobileFilters ? 'hidden md:block' : 'block'}">

      <!-- 筛选器 -->
      <div class="bg-white rounded-2xl shadow-lg p-6 sticky top-6">
        <h3 class="text-lg font-bold text-slate-800 mb-4">筛选器</h3>
        
        <!-- 内容类型 -->
        <div class="mb-6">
          <div class="text-sm font-semibold text-slate-700 mb-3 flex items-center gap-2">
            <span class="w-1 h-4 bg-gradient-to-b from-blue-500 to-purple-500 rounded-full"></span>
            内容类型
          </div>
          <div class="space-y-3">
            <label class="flex items-center gap-3 cursor-pointer group p-2 rounded-lg hover:bg-slate-50 transition-all duration-200">
              <div class="relative">
                <input type="checkbox" bind:checked={filters.types.group} class="sr-only" />
                <div class="w-5 h-5 border-2 border-slate-300 rounded-md flex items-center justify-center transition-all duration-200 group-hover:border-blue-400 {filters.types.group ? 'bg-gradient-to-br from-blue-500 to-blue-600 border-blue-500' : 'bg-white'}">
                  {#if filters.types.group}
                    <svg class="w-3 h-3 text-white" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"/>
                    </svg>
                  {/if}
                </div>
              </div>
              <span class="text-sm text-slate-700 font-medium group-hover:text-slate-900 transition-colors">📱 群组</span>
            </label>
            <label class="flex items-center gap-3 cursor-pointer group p-2 rounded-lg hover:bg-slate-50 transition-all duration-200">
              <div class="relative">
                <input type="checkbox" bind:checked={filters.types.channel} class="sr-only" />
                <div class="w-5 h-5 border-2 border-slate-300 rounded-md flex items-center justify-center transition-all duration-200 group-hover:border-green-400 {filters.types.channel ? 'bg-gradient-to-br from-green-500 to-green-600 border-green-500' : 'bg-white'}">
                  {#if filters.types.channel}
                    <svg class="w-3 h-3 text-white" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"/>
                    </svg>
                  {/if}
                </div>
              </div>
              <span class="text-sm text-slate-700 font-medium group-hover:text-slate-900 transition-colors">📢 频道</span>
            </label>
            <label class="flex items-center gap-3 cursor-pointer group p-2 rounded-lg hover:bg-slate-50 transition-all duration-200">
              <div class="relative">
                <input type="checkbox" bind:checked={filters.types.bot} class="sr-only" />
                <div class="w-5 h-5 border-2 border-slate-300 rounded-md flex items-center justify-center transition-all duration-200 group-hover:border-purple-400 {filters.types.bot ? 'bg-gradient-to-br from-purple-500 to-purple-600 border-purple-500' : 'bg-white'}">
                  {#if filters.types.bot}
                    <svg class="w-3 h-3 text-white" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"/>
                    </svg>
                  {/if}
                </div>
              </div>
              <span class="text-sm text-slate-700 font-medium group-hover:text-slate-900 transition-colors">🤖 机器人</span>
            </label>
            <label class="flex items-center gap-3 cursor-pointer group p-2 rounded-lg hover:bg-slate-50 transition-all duration-200">
              <div class="relative">
                <input type="checkbox" bind:checked={filters.types.message} class="sr-only" />
                <div class="w-5 h-5 border-2 border-slate-300 rounded-md flex items-center justify-center transition-all duration-200 group-hover:border-orange-400 {filters.types.message ? 'bg-gradient-to-br from-orange-500 to-orange-600 border-orange-500' : 'bg-white'}">
                  {#if filters.types.message}
                    <svg class="w-3 h-3 text-white" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"/>
                    </svg>
                  {/if}
                </div>
              </div>
              <span class="text-sm text-slate-700 font-medium group-hover:text-slate-900 transition-colors">💬 消息</span>
            </label>
          </div>
        </div>

        <!-- 时间筛选 -->
        <div class="mb-6">
          <div class="text-sm font-semibold text-slate-700 mb-3 flex items-center gap-2">
            <span class="w-1 h-4 bg-gradient-to-b from-purple-500 to-pink-500 rounded-full"></span>
            时间筛选
          </div>
          <label class="block mb-2">
            <span class="block text-sm font-medium text-slate-700 mb-2 flex items-center gap-1">
              <svg class="w-4 h-4 text-slate-500" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clip-rule="evenodd"/>
              </svg>
              预设时间
            </span>
            <div class="relative">
              <select bind:value={filters.timePreset} class="w-full px-4 py-3 border-2 border-slate-200 rounded-xl bg-white focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all duration-200 appearance-none cursor-pointer hover:border-slate-300 hover:shadow-sm text-sm font-medium text-slate-700">
                <option value="any">🌐 不限</option>
                <option value="1h">⚡ 最近1小时</option>
                <option value="24h">📅 最近24小时</option>
                <option value="7d">📊 最近7天</option>
                <option value="30d">📈 最近30天</option>
                <option value="custom">⚙️ 自定义</option>
              </select>
              <div class="absolute inset-y-0 right-0 flex items-center pr-3 pointer-events-none">
                <svg class="w-5 h-5 text-slate-400" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clip-rule="evenodd"/>
                </svg>
              </div>
            </div>
          </label>
          {#if filters.timePreset === 'custom'}
            <div class="mt-4 p-4 bg-gradient-to-r from-blue-50 to-purple-50 rounded-xl border border-blue-100">
              <div class="flex items-center gap-2 mb-3">
                <svg class="w-4 h-4 text-blue-500" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M6 2a1 1 0 00-1 1v1H4a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V6a2 2 0 00-2-2h-1V3a1 1 0 10-2 0v1H7V3a1 1 0 00-1-1zm0 5a1 1 0 000 2h8a1 1 0 100-2H6z" clip-rule="evenodd"/>
                </svg>
                <span class="text-sm font-semibold text-blue-700">自定义时间范围</span>
              </div>
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
                <label class="block">
                  <span class="block text-xs font-medium text-slate-600 mb-1 flex items-center gap-1">
                    <span class="w-2 h-2 bg-green-400 rounded-full"></span>
                    开始日期
                  </span>
                  <input type="date" bind:value={filters.customStart} class="w-full px-3 py-2.5 border-2 border-slate-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all duration-200 bg-white hover:border-slate-300 text-sm" />
                </label>
                <label class="block">
                  <span class="block text-xs font-medium text-slate-600 mb-1 flex items-center gap-1">
                    <span class="w-2 h-2 bg-red-400 rounded-full"></span>
                    结束日期
                  </span>
                  <input type="date" bind:value={filters.customEnd} class="w-full px-3 py-2.5 border-2 border-slate-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all duration-200 bg-white hover:border-slate-300 text-sm" />
                </label>
              </div>
            </div>
          {/if}
        </div>

        <div class="flex gap-3 mt-6">
          <button class="flex-1 px-4 py-3 text-sm font-semibold text-slate-600 bg-gradient-to-r from-slate-100 to-slate-200 rounded-xl hover:from-slate-200 hover:to-slate-300 transition-all duration-200 shadow-sm hover:shadow-md flex items-center justify-center gap-2" on:click={resetFilters}>
            <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M4 2a1 1 0 011 1v2.101a7.002 7.002 0 0111.601 2.566 1 1 0 11-1.885.666A5.002 5.002 0 005.999 7H9a1 1 0 010 2H4a1 1 0 01-1-1V3a1 1 0 011-1zm.008 9.057a1 1 0 011.276.61A5.002 5.002 0 0014.001 13H11a1 1 0 110-2h5a1 1 0 011 1v5a1 1 0 11-2 0v-2.101a7.002 7.002 0 01-11.601-2.566 1 1 0 01.61-1.276z" clip-rule="evenodd"/>
            </svg>
            重置
          </button>
          <button class="flex-1 px-4 py-3 text-sm font-semibold text-white bg-gradient-to-r from-blue-600 to-blue-700 rounded-xl hover:from-blue-700 hover:to-blue-800 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 shadow-lg hover:shadow-xl flex items-center justify-center gap-2" disabled={isSearching} on:click={() => (page=1, search())}>
            {#if isSearching}
              <svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="m4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              应用中...
            {:else}
              <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M3 3a1 1 0 000 2v8a2 2 0 002 2h2.586l-1.293 1.293a1 1 0 101.414 1.414L10 15.414l2.293 2.293a1 1 0 001.414-1.414L12.414 15H15a2 2 0 002-2V5a1 1 0 100-2H3zm11.707 4.707a1 1 0 00-1.414-1.414L10 9.586 8.707 8.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/>
              </svg>
              应用
            {/if}
          </button>
        </div>
      </div>
    </aside>

    <!-- 右侧结果区域 -->
    <main class="flex-1 min-w-0">
      <!-- 高级选项（折叠） -->
      {#if showAdvanced}
        <section class="mb-6">
          <div class="bg-white rounded-2xl shadow-lg p-6">
            <!-- 来源过滤 -->
            <div class="mb-6">
              <h4 class="text-lg font-semibold text-slate-800 mb-4">来源过滤</h4>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <label class="block">
                  <span class="block text-sm font-medium text-slate-700 mb-1">群组</span>
                  <input placeholder="群组名称" bind:value={filters.sourceGroup} class="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-colors" />
                </label>
                <label class="block">
                  <span class="block text-sm font-medium text-slate-700 mb-1">频道</span>
                  <input placeholder="频道名称" bind:value={filters.sourceChannel} class="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-colors" />
                </label>
                <label class="block">
                  <span class="block text-sm font-medium text-slate-700 mb-1">用户</span>
                  <input placeholder="发送者" bind:value={filters.sourceUser} class="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-colors" />
                </label>
                <label class="block">
                  <span class="block text-sm font-medium text-slate-700 mb-1">语言</span>
                  <input placeholder="如: zh, en" bind:value={filters.language} class="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-colors" />
                </label>
              </div>
            </div>
            <!-- 质量过滤 -->
            <div class="mb-6">
              <h4 class="text-lg font-semibold text-slate-800 mb-4">质量过滤</h4>
              <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                <label class="block">
                  <span class="block text-sm font-medium text-slate-700 mb-1">最小长度</span>
                  <input type="number" min="0" placeholder="字符数" bind:value={filters.minLength} class="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-colors" />
                </label>
                <label class="block">
                  <span class="block text-sm font-medium text-slate-700 mb-1">最小互动</span>
                  <input type="number" min="0" placeholder="点赞/转发数" bind:value={filters.minInteractions} class="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-colors" />
                </label>
                <label class="block">
                  <span class="block text-sm font-medium text-slate-700 mb-1">最小相关性</span>
                  <input type="number" min="0" max="1" step="0.01" placeholder="0~1" bind:value={filters.minRelevance} class="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-colors" />
                </label>
              </div>
            </div>
            <!-- 排序与分页 -->
            <div>
              <h4 class="text-lg font-semibold text-slate-800 mb-4">排序与分页</h4>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <label class="block">
                  <span class="block text-sm font-medium text-slate-700 mb-1">排序</span>
                  <select bind:value={sort} class="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-colors">
                    <option value="relevance">相关性</option>
                    <option value="time">时间</option>
                    <option value="heat">热度</option>
                  </select>
                </label>
                <label class="block">
                  <span class="block text-sm font-medium text-slate-700 mb-1">每页条数</span>
                  <select bind:value={size} class="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-colors">
                    <option value="10">10</option>
                    <option value="20">20</option>
                    <option value="50">50</option>
                  </select>
                </label>
              </div>
            </div>
          </div>
        </section>
      {/if}

      <!-- 搜索结果 -->
      <section>
        <div class="flex flex-wrap sm:flex-nowrap justify-between items-center gap-4 mb-6 bg-white rounded-2xl shadow-lg p-4">
          <div class="flex items-center gap-4 flex-shrink-0">
            <span class="text-slate-700">共 <strong class="text-slate-900">{total}</strong> 条结果</span>
            {#if elapsedMs > 0}
              <span class="text-slate-500 text-sm">耗时 {elapsedMs}ms</span>
            {/if}

          </div>
          <div class="w-full sm:w-auto">
            <div class="flex items-center justify-center sm:justify-end gap-3">
              <button class="flex-shrink-0 px-4 py-2 text-sm font-medium text-slate-600 bg-slate-100 rounded-lg hover:bg-slate-200 disabled:opacity-50 disabled:cursor-not-allowed transition-colors" disabled={isSearching || page <= 1} on:click={() => goToPage(page - 1)}>上一页</button>
              <span class="flex-shrink-0 text-sm text-slate-600 px-3">第 {page} / {totalPages()} 页</span>
              <button class="flex-shrink-0 px-4 py-2 text-sm font-medium text-slate-600 bg-slate-100 rounded-lg hover:bg-slate-200 disabled:opacity-50 disabled:cursor-not-allowed transition-colors" disabled={isSearching || page >= totalPages()} on:click={() => goToPage(page + 1)}>下一页</button>
            </div>
          </div>
        </div>



        <div class="space-y-4 results-container">
          {#if results.length === 0 && !isSearching}
            <div class="bg-white rounded-2xl shadow-lg p-12 text-center">
              <div class="text-slate-400 text-lg">暂无数据，请输入关键字进行搜索。</div>
            </div>
          {/if}
          {#each results as item}
            <article class="bg-white rounded-2xl shadow-lg p-6 hover:shadow-xl transition-shadow duration-300 flex flex-col h-full">
              <div class="flex-grow">
                <a href={item.link} target="_blank" rel="noopener noreferrer" class="text-lg font-semibold text-blue-600 hover:underline">
                  {@html highlight(item.title, query)}
                </a>
                {#if item.description}
                <div class="text-slate-600 leading-relaxed my-3 text-sm">
                  {@html highlight(item.description, query)}
                </div>
                {/if}
              </div>

              <div class="flex items-center text-sm text-slate-500 mt-auto pt-4 border-t border-slate-100">
                {#if item.type === 'group'}
                  <span class="inline-flex items-center">
                    <svg class="w-4 h-4 mr-1.5 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.653-.124-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.653.124-1.283.356-1.857m0 0a3.002 3.002 0 012.288-2.542M11 11a4 4 0 11-8 0 4 4 0 018 0z"></path></svg>
                    Group
                  </span>
                {:else if item.type === 'channel'}
                  <span class="inline-flex items-center">
                    <svg class="w-4 h-4 mr-1.5 text-yellow-500" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.536 8.464a5 5 0 010 7.072m2.828-9.9a9 9 0 010 12.728M5.636 18.364a9 9 0 010-12.728m2.828 9.9a5 5 0 010-7.072"></path></svg>
                    Channel
                  </span>
                {/if}
                <span class="mx-2">·</span>
                <span class="inline-flex items-center">
                  <svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M15 21v-1a6 6 0 00-5.197-5.92M9 21v-1a6 6 0 016-6"></path></svg>
                  {item.members_count} members
                </span>

                <div class="flex-grow"></div>

                {#if item.type === 'group' || item.type === 'channel'}
                    <button class="bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded-lg text-sm transition-all duration-300 shadow-md hover:shadow-lg" on:click={() => openGroupModal(item)} title="查看群组统计详情">查看详情</button>
                {:else}
                    <button class="text-slate-400 text-sm font-medium cursor-not-allowed" disabled title="占位">查看详情</button>
                {/if}
              </div>
            </article>
          {/each}
        </div>
      </section>
    </main>
  </div>
</div>

<!-- 群组详情弹窗组件 -->
<GroupModal 
  show={showGroupModal} 
  groupData={selectedGroupData} 
  on:close={closeGroupModal} 
/>

<!-- 收录弹窗组件 -->
<CollectionModal 
  bind:show={showCollectionModal} 
  on:close={() => showCollectionModal = false}
/>



<!--
  关键算法说明：
  - highlight(text, q): 简单分词高亮（按空格拆分，忽略大小写），避免XSS的转义处理
  - debounce(fn, wait): 输入建议请求防抖
  待优化事项：
  - 接入真实鉴权与环境变量；完善机器人状态、来源跳转、收藏/分享；分页组件化；滚动加载可选
  - 搜索建议支持键盘方向键选择；高亮匹配片段截断与上下文偏移
  兼容性说明：
  - 客户端渲染，localStorage 仅在浏览器环境访问；样式兼容现代浏览器
-->

<footer class="text-center py-4 text-gray-500 text-sm">
  <span>© 2024 Telegram搜索平台 - 高效搜索, 精准定位</span>
  <a href="#" on:click|preventDefault={openDisclaimerModal} class="hover:underline text-red-500 ml-2">免责声明</a>
</footer>

<DisclaimerModal bind:show={showDisclaimerModal} on:close={closeDisclaimerModal} />
