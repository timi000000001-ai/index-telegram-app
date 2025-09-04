<!--
 * Telegram搜索平台 - 全局布局组件
 * @description 提供全站通用的布局结构和样式
 * @author 前端工程师
 * @date 2024-01-23
 * @version 1.0.0
 * © Telegram搜索平台
-->

<script>
  /**
   * 全局布局逻辑
   * @description 处理全站通用的状态管理和初始化
   */
  import '../app.css';
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { base } from '$app/paths';
  
  let user = null;
  let isMobileMenuOpen = false;

  function toggleMobileMenu() {
    isMobileMenuOpen = !isMobileMenuOpen;
  }

  /**
   * 当前页面路径
   * @description 用于导航菜单的选中状态判断
   */
  $: currentPath = $page.url.pathname;
  
  /**
   * 页面初始化
   * @description 设置全局配置和主题
   */
  onMount(() => {
    // 设置页面标题
    document.title = 'Telegram搜索平台';
    
    // 设置favicon
    /** @type {HTMLLinkElement} */
    const favicon = document.querySelector('link[rel="icon"]') || document.createElement('link');
    favicon.rel = 'icon';
    favicon.href = '/favicon.svg';
    favicon.type = 'image/svg+xml';
    if (!document.querySelector('link[rel="icon"]')) {
      document.head.appendChild(favicon);
    }

    if (window.Telegram && window.Telegram.WebApp) {
      const tg = window.Telegram.WebApp;
      if (tg.initDataUnsafe && tg.initDataUnsafe.user) {
        user = tg.initDataUnsafe.user;
      }
    }
  });
</script>

<!--
 * 全局HTML结构
 * @description 提供响应式布局容器
-->
<div class="min-h-screen bg-gradient-to-br from-slate-50 to-blue-50">
  <!-- 头部导航栏 -->
  <header class="bg-white/80 backdrop-blur-md border-b border-slate-200/50 sticky top-0 z-50">
    <nav class="container mx-auto px-4 py-4">
      <div class="flex items-center justify-between">
        <!-- Logo和网站名称 -->
        <div class="flex items-center gap-3">
          <div class="w-8 h-8 bg-gradient-to-br from-blue-500 to-purple-600 rounded-lg flex items-center justify-center">
            <svg class="w-5 h-5 text-white" fill="currentColor" viewBox="0 0 24 24">
              <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
            </svg>
          </div>
          <a href="{base}/" class="text-xl font-bold text-slate-800 hover:text-blue-600 transition-colors">
            Telegram搜索
          </a>
        </div>
        
        <!-- 导航菜单 -->
         <div class="hidden md:flex items-center gap-6">
           <a href="{base}/" 
              class="nav-link {currentPath === `${base}/` ? 'nav-active' : ''}"
              data-text="首页">
             首页
           </a>
           <a href="{base}/daily-stats" 
              class="nav-link {currentPath === `${base}/daily-stats` ? 'nav-active' : ''}"
              data-text="每日动态">
             每日动态
           </a>
           <a href="{base}/daily-new" 
              class="nav-link {currentPath === `${base}/daily-new` ? 'nav-active' : ''}"
              data-text="每日新增">
             每日新增
           </a>
           <a href="{base}/bots" 
              class="nav-link {currentPath === `${base}/bots` ? 'nav-active' : ''}"
              data-text="机器人管理">
             机器人管理
           </a>
         </div>
        
        <!-- 登录按钮和用户菜单 -->
        <div class="flex items-center gap-3">
          <!-- 移动端菜单按钮 -->
          <button on:click={toggleMobileMenu} class="md:hidden p-2 text-slate-600 hover:text-blue-600 transition-colors" aria-label="打开移动端菜单">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"></path>
            </svg>
          </button>
          
          <!-- 登录按钮 -->
          {#if user}
            <span class="text-slate-800 font-medium">{user.first_name} {user.last_name || ''}</span>
          {:else}
            <a href="{base}/login" class="bg-gradient-to-r from-blue-500 to-purple-600 text-white px-4 py-2 rounded-lg font-medium hover:shadow-lg transition-all duration-300">
              登录
            </a>
          {/if}
        </div>
      </div>
    </nav>

    <!-- 移动端菜单 -->
    {#if isMobileMenuOpen}
      <div class="absolute left-0 right-0 bg-white shadow-lg md:hidden px-4 pt-2 pb-4 space-y-2 z-50">
        <a href="{base}/" on:click={() => isMobileMenuOpen = false} class="block nav-link {currentPath === `${base}/` ? 'nav-active' : ''}">首页</a>
        <a href="{base}/daily-stats" on:click={() => isMobileMenuOpen = false} class="block nav-link {currentPath === `${base}/daily-stats` ? 'nav-active' : ''}">每日动态</a>
        <a href="{base}/daily-new" on:click={() => isMobileMenuOpen = false} class="block nav-link {currentPath === `${base}/daily-new` ? 'nav-active' : ''}">每日新增</a>
        <a href="{base}/bots" on:click={() => isMobileMenuOpen = false} class="block nav-link {currentPath === `${base}/bots` ? 'nav-active' : ''}">机器人管理</a>
      </div>
    {/if}
  </header>
  
  <!-- 主内容区域 -->
  <main class="container mx-auto px-4 py-6">
    <slot />
  </main>
  

</div>

<!--
 * 样式定义
 * @description 全局样式覆盖和自定义
-->
<style>
  :global(body) {
    font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    line-height: 1.6;
  }
  
  :global(*) {
    box-sizing: border-box;
  }
  
  :global(html) {
    scroll-behavior: smooth;
  }
  
  /* 导航链接样式 */
  .nav-link {
    @apply relative px-3 py-2 text-slate-600 font-medium transition-all duration-300 ease-out;
    @apply hover:text-blue-600;
  }
  
  /* 选中状态样式 */
  .nav-active {
    @apply text-blue-600;
  }
  
  /* 底部指示器动画 */
  .nav-link::after {
    content: '';
    @apply absolute bottom-0 left-1/2 w-0 h-0.5 bg-gradient-to-r from-blue-500 to-purple-600;
    @apply transition-all duration-300 ease-out;
    transform: translateX(-50%);
  }
  
  .nav-link:hover::after {
    @apply w-full;
  }
  
  .nav-active::after {
    @apply w-full;
  }
  
  /* 点击波纹效果 */
  .nav-link {
    overflow: hidden;
  }
  
  .nav-link::before {
    content: '';
    @apply absolute inset-0 bg-blue-100 rounded-lg opacity-0;
    @apply transition-all duration-200 ease-out;
    transform: scale(0.8);
  }
  
  .nav-link:active::before {
    @apply opacity-30;
    transform: scale(1);
  }
  
  /* 文字渐变效果（选中状态） */
  .nav-active {
    background: linear-gradient(135deg, #3b82f6, #8b5cf6);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }
</style>

<!--
 * 文件说明
 * @description 全局布局组件，为所有页面提供统一的结构和样式
 * @features 响应式设计、主题支持、SEO优化
 * @compatibility 支持现代浏览器，兼容移动端
-->