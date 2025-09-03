<script>
/**
 * @fileoverview 收录弹窗组件
 * @description 支持单个和多个链接收录，以及通过TG机器人收录
 * @author 前端工程师
 * @date 2024-01-20
 * @version 1.0.0
 * © Telegram Search Platform
 */

import { createEventDispatcher } from 'svelte';

const dispatch = createEventDispatcher();

/** @type {boolean} 是否显示弹窗 */
export let show = false;

/** @type {string} 收录模式：'single' | 'multiple' | 'telegram' */
let mode = 'single';

/** @type {string} 单个链接输入 */
let singleLink = '';

/** @type {string} 多个链接输入（每行一个） */
let multipleLinks = '';

/** @type {string} 链接标题 */
let linkTitle = '';

/** @type {string} 链接描述 */
let linkDescription = '';

/** @type {string} 分类标签 */
let category = '';

/** @type {boolean} 是否正在保存 */
let isSaving = false;

/** @type {string} 提示消息 */
let message = '';

/** @type {'success' | 'error' | ''} 消息类型 */
let messageType = '';

/** @type {string} 链接验证错误信息 */
let linkError = '';

/** @type {string} 批量链接验证错误信息 */
let multipleLinksError = '';

/**
 * 关闭弹窗
 */
function closeModal() {
  show = false;
  resetForm();
  dispatch('close');
}

/**
 * 重置表单
 */
function resetForm() {
  singleLink = '';
  multipleLinks = '';
  linkTitle = '';
  linkDescription = '';
  category = '';
  message = '';
  messageType = '';
  linkError = '';
  multipleLinksError = '';
}

/**
 * 验证链接格式
 * @param {string} link - 要验证的链接
 * @returns {boolean} 是否为有效的Telegram链接
 */
function isValidTelegramLink(link) {
  return link.trim().startsWith('https://t.me');
}

/**
 * 验证单个链接
 */
function validateSingleLink() {
  if (!singleLink.trim()) {
    linkError = '';
    return true;
  }
  
  if (!isValidTelegramLink(singleLink)) {
    linkError = '只能添加 https://t.me 开头的链接';
    return false;
  }
  
  linkError = '';
  return true;
}

/**
 * 验证批量链接
 */
function validateMultipleLinks() {
  if (!multipleLinks.trim()) {
    multipleLinksError = '';
    return true;
  }
  
  const links = multipleLinks.split('\n').filter(link => link.trim());
  const invalidLinks = links.filter(link => !isValidTelegramLink(link));
  
  if (invalidLinks.length > 0) {
    multipleLinksError = `只能添加 https://t.me 开头的链接，发现 ${invalidLinks.length} 个无效链接`;
    return false;
  }
  
  multipleLinksError = '';
  return true;
}

/**
 * 保存收录
 * @param {boolean} continueAfter - 保存后是否继续收录
 */
async function saveCollection(continueAfter = false) {
  if (isSaving) return;
  
  // 验证输入
  if (mode === 'single') {
    if (!singleLink.trim()) {
      showMessage('请输入链接地址', 'error');
      return;
    }
    if (!validateSingleLink()) {
      return;
    }
  }
  
  if (mode === 'multiple') {
    if (!multipleLinks.trim()) {
      showMessage('请输入链接地址', 'error');
      return;
    }
    if (!validateMultipleLinks()) {
      return;
    }
  }
  
  isSaving = true;
  
  try {
    // 准备收录数据
    const collectionData = {
      mode,
      title: linkTitle,
      description: linkDescription,
      category,
      links: mode === 'single' ? [singleLink] : multipleLinks.split('\n').filter(link => link.trim()),
      timestamp: new Date().toISOString()
    };
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000));
    
    // 随机成功/失败（用于演示）
    const success = Math.random() > 0.2;
    
    if (success) {
      showMessage('收录成功！', 'success');
      
      if (continueAfter) {
        // 清空链接输入，保留其他信息
        singleLink = '';
        multipleLinks = '';
      } else {
        // 延迟关闭弹窗
        setTimeout(() => {
          closeModal();
        }, 1500);
      }
    } else {
      showMessage('收录失败，请重试', 'error');
    }
    
  } catch (error) {
    console.error('收录失败:', error);
    showMessage('收录失败，请重试', 'error');
  } finally {
    isSaving = false;
  }
}

/**
 * 显示消息
 * @param {string} msg - 消息内容
 * @param {'success' | 'error'} type - 消息类型
 */
function showMessage(msg, type) {
  message = msg;
  messageType = type;
  
  // 3秒后清除消息
  setTimeout(() => {
    message = '';
    messageType = '';
  }, 3000);
}

/**
 * 切换收录模式
 * @param {string} newMode - 新模式
 */
function switchMode(newMode) {
  mode = newMode;
  message = '';
  messageType = '';
  linkError = '';
  multipleLinksError = '';
}

/**
 * 打开Telegram机器人
 */
function openTelegramBot() {
  const botUrl = 'https://t.me/your_collection_bot';
  window.open(botUrl, '_blank');
  showMessage('请在Telegram中向机器人发送链接进行收录', 'success');
}
</script>

<!-- 弹窗遮罩 -->
{#if show}
  <div class="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center p-4" 
       role="dialog" 
       aria-modal="true" 
       tabindex="-1"
       on:click={closeModal} 
       on:keydown={(e) => e.key === 'Escape' && closeModal()}>
    <!-- 弹窗内容 -->
    <div class="bg-white rounded-2xl shadow-2xl w-full max-w-2xl max-h-[90vh] overflow-hidden"
         role="document"
         tabindex="0"
         on:click|stopPropagation
         on:keydown|stopPropagation>
      <!-- 弹窗头部 -->
      <div class="bg-gradient-to-r from-blue-500 to-purple-600 text-white p-6">
        <div class="flex items-center justify-between">
          <h2 class="text-xl font-bold">添加收录</h2>
          <button class="text-white/80 hover:text-white transition-colors" 
                  on:click={closeModal} 
                  aria-label="关闭弹窗">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
          </button>
        </div>
        
        <!-- 模式切换 -->
        <div class="flex gap-2 mt-4">
          <button 
            class="px-4 py-2 rounded-lg text-sm font-medium transition-all {mode === 'single' ? 'bg-white/20 text-white' : 'bg-white/10 text-white/70 hover:bg-white/15'}"
            on:click={() => switchMode('single')}
          >
            单个链接
          </button>
          <button 
            class="px-4 py-2 rounded-lg text-sm font-medium transition-all {mode === 'multiple' ? 'bg-white/20 text-white' : 'bg-white/10 text-white/70 hover:bg-white/15'}"
            on:click={() => switchMode('multiple')}
          >
            批量链接
          </button>
          <button 
            class="px-4 py-2 rounded-lg text-sm font-medium transition-all {mode === 'telegram' ? 'bg-white/20 text-white' : 'bg-white/10 text-white/70 hover:bg-white/15'}"
            on:click={() => switchMode('telegram')}
          >
            TG机器人
          </button>
        </div>
      </div>
      
      <!-- 弹窗内容区域 -->
      <div class="p-6 overflow-y-auto max-h-[60vh]">
        {#if mode === 'single'}
          <!-- 单个链接收录 -->
          <div class="space-y-4">
            <div>
              <label for="single-link" class="block text-sm font-medium text-gray-700 mb-2">链接地址 *</label>
              <input 
                id="single-link"
                type="url" 
                class="w-full px-4 py-3 border rounded-lg outline-none transition-colors {linkError ? 'border-red-500 focus:ring-2 focus:ring-red-500 focus:border-red-500' : 'border-gray-300 focus:ring-2 focus:ring-blue-500 focus:border-blue-500'}"
                placeholder="https://t.me/example"
                bind:value={singleLink}
                on:input={validateSingleLink}
                on:blur={validateSingleLink}
              />
              {#if linkError}
                <p class="text-red-500 text-sm mt-1 flex items-center gap-1">
                  <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
                  </svg>
                  {linkError}
                </p>
              {/if}
            </div>
            
            <div>
              <label for="link-title" class="block text-sm font-medium text-gray-700 mb-2">标题</label>
              <input 
                id="link-title"
                type="text" 
                class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-colors"
                placeholder="链接标题（可选）"
                bind:value={linkTitle}
              />
            </div>
            
            <div>
              <label for="link-description" class="block text-sm font-medium text-gray-700 mb-2">描述</label>
              <textarea 
                id="link-description"
                class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-colors resize-none"
                rows="3"
                placeholder="链接描述（可选）"
                bind:value={linkDescription}
              ></textarea>
            </div>
            
            <div>
              <label for="category" class="block text-sm font-medium text-gray-700 mb-2">分类</label>
              <select 
                id="category"
                class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-colors"
                bind:value={category}
              >
                <option value="">选择分类（可选）</option>
                <option value="technology">技术</option>
                <option value="news">新闻</option>
                <option value="entertainment">娱乐</option>
                <option value="education">教育</option>
                <option value="business">商业</option>
                <option value="other">其他</option>
              </select>
            </div>
          </div>
          
        {:else if mode === 'multiple'}
          <!-- 批量链接收录 -->
          <div class="space-y-4">
            <div>
              <label for="multiple-links" class="block text-sm font-medium text-gray-700 mb-2">链接地址 *</label>
              <textarea 
                id="multiple-links"
                class="w-full px-4 py-3 border rounded-lg outline-none transition-colors resize-none {multipleLinksError ? 'border-red-500 focus:ring-2 focus:ring-red-500 focus:border-red-500' : 'border-gray-300 focus:ring-2 focus:ring-blue-500 focus:border-blue-500'}"
                rows="6"
                placeholder="每行输入一个链接地址：&#10;https://t.me/example1&#10;https://t.me/example2&#10;https://t.me/example3"
                bind:value={multipleLinks}
                on:input={validateMultipleLinks}
                on:blur={validateMultipleLinks}
              ></textarea>
              {#if multipleLinksError}
                <p class="text-red-500 text-sm mt-1 flex items-center gap-1">
                  <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
                  </svg>
                  {multipleLinksError}
                </p>
              {:else}
                <p class="text-sm text-gray-500 mt-1">每行输入一个Telegram链接，支持批量收录</p>
              {/if}
            </div>
            
            <div>
              <label for="batch-category" class="block text-sm font-medium text-gray-700 mb-2">统一分类</label>
              <select 
                id="batch-category"
                class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-colors"
                bind:value={category}
              >
                <option value="">选择分类（可选）</option>
                <option value="technology">技术</option>
                <option value="news">新闻</option>
                <option value="entertainment">娱乐</option>
                <option value="education">教育</option>
                <option value="business">商业</option>
                <option value="other">其他</option>
              </select>
            </div>
          </div>
          
        {:else if mode === 'telegram'}
          <!-- Telegram机器人收录 -->
          <div class="text-center space-y-6">
            <div class="bg-blue-50 rounded-xl p-6">
              <div class="text-6xl mb-4">🤖</div>
              <h3 class="text-lg font-semibold text-gray-800 mb-2">通过Telegram机器人收录</h3>
              <p class="text-gray-600 mb-4">点击下方按钮打开Telegram机器人，直接发送链接即可快速收录</p>
              
              <button 
                class="bg-blue-500 hover:bg-blue-600 text-white px-6 py-3 rounded-lg font-medium transition-colors inline-flex items-center gap-2"
                on:click={openTelegramBot}
              >
                <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M12 0C5.374 0 0 5.373 0 12s5.374 12 12 12 12-5.373 12-12S18.626 0 12 0zm5.568 8.16c-.169 1.858-.896 6.728-.896 6.728-.377 2.617-1.407 3.08-2.896 1.596l-2.123-1.596-1.018.684c-.215.145-.395.26-.81.26-.528 0-.434-.2-.612-.706L8.4 13.116l-2.917-.93c-.631-.2-.636-.63.14-.936l11.395-4.29c.523-.2 1.01.12.55 1.2z"/>
                </svg>
                打开收录机器人
              </button>
            </div>
            
            <div class="bg-gray-50 rounded-xl p-4">
              <h4 class="font-medium text-gray-800 mb-2">使用说明：</h4>
              <ul class="text-sm text-gray-600 space-y-1 text-left">
                <li>• 点击按钮打开Telegram机器人</li>
                <li>• 直接发送链接给机器人</li>
                <li>• 机器人会自动解析并收录链接</li>
                <li>• 支持批量发送多个链接</li>
              </ul>
            </div>
          </div>
        {/if}
        
        <!-- 消息提示 -->
        {#if message}
          <div class="mt-4 p-3 rounded-lg {messageType === 'success' ? 'bg-green-50 text-green-700 border border-green-200' : 'bg-red-50 text-red-700 border border-red-200'}">
            <div class="flex items-center gap-2">
              {#if messageType === 'success'}
                <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"></path>
                </svg>
              {:else}
                <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd"></path>
                </svg>
              {/if}
              <span>{message}</span>
            </div>
          </div>
        {/if}
      </div>
      
      <!-- 弹窗底部按钮 -->
      {#if mode !== 'telegram'}
        <div class="bg-gray-50 px-6 py-4 flex gap-3 justify-end">
          <button 
            class="px-6 py-2 text-gray-600 hover:text-gray-800 transition-colors"
            on:click={closeModal}
            disabled={isSaving}
          >
            取消
          </button>
          
          <button 
            class="px-6 py-2 bg-blue-500 hover:bg-blue-600 text-white rounded-lg font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            on:click={() => saveCollection(true)}
            disabled={isSaving}
          >
            {#if isSaving}保存中...{:else}保存后继续收录{/if}
          </button>
          
          <button 
            class="px-6 py-2 bg-green-500 hover:bg-green-600 text-white rounded-lg font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            on:click={() => saveCollection(false)}
            disabled={isSaving}
          >
            {#if isSaving}保存中...{:else}保存{/if}
          </button>
        </div>
      {:else}
        <div class="bg-gray-50 px-6 py-4 flex justify-end">
          <button 
            class="px-6 py-2 text-gray-600 hover:text-gray-800 transition-colors"
            on:click={closeModal}
          >
            关闭
          </button>
        </div>
     
      {/if}
    </div>
  </div>   
{/if}