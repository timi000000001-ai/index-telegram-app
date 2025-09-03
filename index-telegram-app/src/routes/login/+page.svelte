<!--
  文件功能: Telegram 登录与账号管理页面
  主要组件: LoginPage（提供 API 凭据保存、手机号登录、验证码提交、账号列表与采集触发）
  修改历史:
  - 2025-09-02 v0.1.0 初始从首页拆分
  作者: Trae AI
  @date 2025-09-02
  @version 0.1.0
  © 公司/组织
-->
<script>
  import { writable } from 'svelte/store';
  import { browser } from '$app/environment';
  import { onMount } from 'svelte';
  import { apiFetch } from '$lib/index.js';
  import { error as toastError, success as toastSuccess, info as toastInfo } from '$lib/message.js';

  // 页面状态：API 凭据与登录信息
  let apiId = '';
  let apiHash = '';
  /** @type {import('svelte/store').Writable<Array<{phone: string, session: any}>>} */
  let accounts = writable([]);
  let phoneNumber = '';
  let code = '';
  let codeRequired = false;
  let error = '';
  let isSubmitting = false;

  // 生命周期：仅在浏览器读取/写入本地存储，避免 SSR 报错
  onMount(() => {
    if (browser) {
      apiId = localStorage.getItem('api_id') || '';
      apiHash = localStorage.getItem('api_hash') || '';
      try {
        const storedAccounts = JSON.parse(localStorage.getItem('accounts') || '[]');
        accounts.set(storedAccounts);
      } catch (e) {
        console.error('Failed to parse accounts from localStorage', e);
      }
    }
  });

  /**
   * 保存 API 凭据
   * @returns {void}
   */
  function saveApiCredentials() {
    localStorage.setItem('api_id', apiId);
    localStorage.setItem('api_hash', apiHash);
    toastSuccess('API 凭据已保存');
  }

  /**
   * 发起登录请求（手机号）
   * @returns {Promise<void>}
   */
  async function startLogin() {
    error = '';
    isSubmitting = true;
    try {
      const { ok, data: result, error: errMsg } = await apiFetch('/api/login', {
        method: 'POST',
        body: { api_id: parseInt(apiId), api_hash: apiHash, phone_number: phoneNumber }
      });
      if (!ok) {
        error = errMsg || 'Login failed';
        toastError(error);
      } else if (result?.status === 'code_required') {
        codeRequired = true;
        toastInfo('已发送验证码，请输入');
      } else if (result?.status === 'success') {
        accounts.update((accs) => {
          const newAccounts = [...accs, { phone: phoneNumber, session: result.session }];
          localStorage.setItem('accounts', JSON.stringify(newAccounts));
          return newAccounts;
        });
        phoneNumber = '';
        codeRequired = false;
        toastSuccess('登录成功');
      } else {
        error = (result && result.error) || 'Login failed';
        toastError(error);
      }
    } catch (/** @type {any} */ err) {
      error = 'Network error: ' + err.message;
      toastError(error);
    } finally {
      isSubmitting = false;
    }
  }

  /**
   * 提交验证码
   * @returns {Promise<void>}
   */
  async function submitCode() {
    error = '';
    isSubmitting = true;
    try {
      const { ok, data: result, error: errMsg } = await apiFetch('/api/verify', {
        method: 'POST',
        body: { api_id: parseInt(apiId), api_hash: apiHash, phone_number: phoneNumber, code }
      });
      if (!ok) {
        error = errMsg || 'Verification failed';
        toastError(error);
      } else if (result?.status === 'success') {
        accounts.update((accs) => {
          const newAccounts = [...accs, { phone: phoneNumber, session: result.session }];
          localStorage.setItem('accounts', JSON.stringify(newAccounts));
          return newAccounts;
        });
        phoneNumber = '';
        code = '';
        codeRequired = false;
        toastSuccess('验证成功，已登录');
      } else {
        error = (result && result.error) || 'Verification failed';
        toastError(error);
      }
    } catch (/** @type {any} */ err) {
      error = 'Network error: ' + err.message;
      toastError(error);
    } finally {
      isSubmitting = false;
    }
  }

  /**
   * 触发数据采集
   * @param {string} phone 手机号
   * @returns {Promise<void>}
   */
  async function startCollection(phone) {
    try {
      const { ok, data: result, error: errMsg } = await apiFetch('/api/collect', {
        method: 'POST',
        body: { phone_number: phone }
      });
      if (!ok) {
        toastError(errMsg || '采集启动失败');
      } else if (result?.status === 'success') {
        toastSuccess('已启动采集: ' + phone);
      } else {
        error = (result && result.error) || 'Collection failed';
        toastError(error);
      }
    } catch (/** @type {any} */ err) {
      error = 'Network error: ' + err.message;
      toastError(error);
    }
  }
</script>

<main class="min-h-screen bg-gradient-to-br from-slate-50 via-white to-slate-100 py-8 px-4">
  <div class="max-w-4xl mx-auto">
    <div class="text-center mb-8">
      <h1 class="text-4xl font-bold bg-gradient-to-r from-sky-600 via-blue-600 to-indigo-600 bg-clip-text text-transparent mb-2">
        Telegram Data Collector
      </h1>
      <p class="text-slate-600 text-lg">管理您的 Telegram API 凭据和账号登录</p>
    </div>

    <!-- API Credentials Card -->
    <div class="bg-white rounded-xl shadow-lg border border-slate-200 p-6 mb-6 hover:shadow-xl transition-all duration-300">
      <div class="flex items-center gap-3 mb-6">
        <div class="w-10 h-10 bg-gradient-to-r from-emerald-500 to-teal-500 rounded-lg flex items-center justify-center">
          <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m0 0a2 2 0 012 2m-2-2v6m0 0V9a2 2 0 00-2-2M9 7a2 2 0 00-2 2v6a2 2 0 002 2h6a2 2 0 002-2V9a2 2 0 00-2-2"></path>
          </svg>
        </div>
        <h2 class="text-xl font-semibold text-slate-800">API 凭据配置</h2>
      </div>
      <div class="space-y-4">
        <div>
          <label for="apiId" class="block text-sm font-medium text-slate-700 mb-2">API ID</label>
          <input 
            id="apiId" 
            bind:value={apiId} 
            placeholder="请输入 API ID" 
            type="number" 
            class="w-full px-4 py-3 border border-slate-300 rounded-lg focus:ring-2 focus:ring-sky-500 focus:border-sky-500 transition-colors duration-200 text-slate-900 placeholder-slate-400"
          />
        </div>
        <div>
          <label for="apiHash" class="block text-sm font-medium text-slate-700 mb-2">API Hash</label>
          <input 
            id="apiHash" 
            bind:value={apiHash} 
            placeholder="请输入 API Hash" 
            class="w-full px-4 py-3 border border-slate-300 rounded-lg focus:ring-2 focus:ring-sky-500 focus:border-sky-500 transition-colors duration-200 text-slate-900 placeholder-slate-400"
          />
        </div>
        <button 
          class="w-full bg-gradient-to-r from-emerald-500 to-teal-500 hover:from-emerald-600 hover:to-teal-600 text-white font-medium py-3 px-6 rounded-lg transition-all duration-200 transform hover:scale-[1.02] hover:shadow-lg disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none" 
          aria-busy={isSubmitting} 
          disabled={isSubmitting} 
          on:click={saveApiCredentials}
        >
          {#if isSubmitting}
            <div class="flex items-center justify-center gap-2">
              <div class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
              保存中...
            </div>
          {:else}
            保存凭据
          {/if}
        </button>
      </div>
    </div>

    <!-- Login Card -->
    <div class="bg-white rounded-xl shadow-lg border border-slate-200 p-6 mb-6 hover:shadow-xl transition-all duration-300">
      <div class="flex items-center gap-3 mb-6">
        <div class="w-10 h-10 bg-gradient-to-r from-blue-500 to-indigo-500 rounded-lg flex items-center justify-center">
          <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
          </svg>
        </div>
        <h2 class="text-xl font-semibold text-slate-800">账号登录</h2>
      </div>
      <div class="space-y-4">
        <div>
          <label for="phoneNumber" class="block text-sm font-medium text-slate-700 mb-2">手机号码</label>
          <input 
            id="phoneNumber" 
            bind:value={phoneNumber} 
            placeholder="+1234567890" 
            class="w-full px-4 py-3 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors duration-200 text-slate-900 placeholder-slate-400"
          />
        </div>
        <button 
          class="w-full bg-gradient-to-r from-blue-500 to-indigo-500 hover:from-blue-600 hover:to-indigo-600 text-white font-medium py-3 px-6 rounded-lg transition-all duration-200 transform hover:scale-[1.02] hover:shadow-lg disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none" 
          aria-busy={isSubmitting} 
          disabled={isSubmitting} 
          on:click={startLogin}
        >
          {#if isSubmitting}
            <div class="flex items-center justify-center gap-2">
              <div class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
              登录中...
            </div>
          {:else}
            开始登录
          {/if}
        </button>
        
        {#if codeRequired}
          <div class="mt-6 p-4 bg-amber-50 border border-amber-200 rounded-lg">
            <div class="flex items-center gap-2 mb-3">
              <svg class="w-5 h-5 text-amber-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"></path>
              </svg>
              <span class="text-amber-800 font-medium">验证码验证</span>
            </div>
            <div class="space-y-3">
              <div>
                <label for="code" class="block text-sm font-medium text-amber-700 mb-2">验证码</label>
                <input 
                  id="code" 
                  bind:value={code} 
                  placeholder="请输入验证码" 
                  class="w-full px-4 py-3 border border-amber-300 rounded-lg focus:ring-2 focus:ring-amber-500 focus:border-amber-500 transition-colors duration-200 text-slate-900 placeholder-amber-400 bg-white"
                />
              </div>
              <button 
                class="w-full bg-gradient-to-r from-amber-500 to-orange-500 hover:from-amber-600 hover:to-orange-600 text-white font-medium py-3 px-6 rounded-lg transition-all duration-200 transform hover:scale-[1.02] hover:shadow-lg disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none" 
                aria-busy={isSubmitting} 
                disabled={isSubmitting} 
                on:click={submitCode}
              >
                {#if isSubmitting}
                  <div class="flex items-center justify-center gap-2">
                    <div class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                    验证中...
                  </div>
                {:else}
                  提交验证码
                {/if}
              </button>
            </div>
          </div>
        {/if}
      </div>
    </div>

    <!-- Accounts Card -->
    <div class="bg-white rounded-xl shadow-lg border border-slate-200 p-6 hover:shadow-xl transition-all duration-300">
      <div class="flex items-center gap-3 mb-6">
        <div class="w-10 h-10 bg-gradient-to-r from-purple-500 to-pink-500 rounded-lg flex items-center justify-center">
          <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"></path>
          </svg>
        </div>
        <h2 class="text-xl font-semibold text-slate-800">已登录账号</h2>
      </div>
      
      {#if error}
        <div class="mb-4 p-4 bg-red-50 border border-red-200 rounded-lg">
          <div class="flex items-center gap-2">
            <svg class="w-5 h-5 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
            </svg>
            <p class="text-red-800 font-medium">{error}</p>
          </div>
        </div>
      {/if}
      
      {#if $accounts.length === 0}
        <div class="text-center py-12">
          <div class="w-16 h-16 bg-slate-100 rounded-full flex items-center justify-center mx-auto mb-4">
            <svg class="w-8 h-8 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
            </svg>
          </div>
          <p class="text-slate-500 text-lg">暂无已登录账号</p>
          <p class="text-slate-400 text-sm mt-1">请先配置 API 凭据并登录账号</p>
        </div>
      {:else}
        <div class="space-y-3">
          {#each $accounts as account, index}
            <div class="flex items-center justify-between p-4 bg-slate-50 hover:bg-slate-100 rounded-lg border border-slate-200 transition-colors duration-200">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 bg-gradient-to-r from-green-400 to-blue-500 rounded-full flex items-center justify-center">
                  <span class="text-white font-semibold text-sm">{index + 1}</span>
                </div>
                <div>
                  <p class="font-medium text-slate-800">{account.phone}</p>
                  <p class="text-sm text-slate-500">已连接</p>
                </div>
              </div>
              <button 
                class="bg-gradient-to-r from-purple-500 to-pink-500 hover:from-purple-600 hover:to-pink-600 text-white font-medium py-2 px-4 rounded-lg transition-all duration-200 transform hover:scale-105 hover:shadow-md disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none text-sm" 
                aria-busy={isSubmitting} 
                disabled={isSubmitting} 
                on:click={() => startCollection(account.phone)}
              >
                {#if isSubmitting}
                  <div class="flex items-center gap-2">
                    <div class="w-3 h-3 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                    采集中
                  </div>
                {:else}
                  开始采集
                {/if}
              </button>
            </div>
          {/each}
        </div>
      {/if}
    </div>
  </div>
</main>



<!--
  关键算法说明：本页仅包含前端表单与简单的 fetch 交互，无复杂算法。
  待优化事项：抽离 API 基地址与鉴权为环境变量；增加按钮 Loading 与防抖；统一错误提示组件化。
  兼容性说明：已通过 onMount + browser 判断避免 SSR 下的 localStorage 访问。
-->