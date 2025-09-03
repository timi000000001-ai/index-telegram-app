<script lang="ts">
/**
 * @fileoverview 打字机效果文本动画组件
 * @description 实现文本逐字显示的动画效果，支持HTML内容和自定义播放速度
 * @author 前端工程师
 * @date 2024-01-20
 * @version 1.0.0
 * © Telegram Search Platform
 */

import { onMount, onDestroy } from 'svelte';
import { browser } from '$app/environment';

// ===================== 组件属性 =====================
/** @type {string} 要显示的文本内容（支持HTML） */
export let content = '';
/** @type {number} 打字速度（毫秒/字符） */
export let speed = 50;
/** @type {boolean} 是否自动开始播放 */
export let autoStart = true;
/** @type {boolean} 是否循环播放 */
export let loop = false;
/** @type {boolean} 是否显示光标 */
export let showCursor = true;

// ===================== 状态管理 =====================
/** @type {string} 当前显示的文本 */
let displayedText = '';
/** @type {boolean} 是否正在播放动画 */
let isPlaying = false;
/** @type {boolean} 动画是否已完成 */
let isComplete = false;
/** @type {number} 当前字符索引 */
let currentIndex = 0;
/** @type {number|null} 定时器ID */
let timerId: number | null = null;
/** @type {string} 纯文本内容（去除HTML标签） */
let plainText = '';
/** @type {boolean} 光标是否可见 */
let cursorVisible = true;
/** @type {number|null} 光标闪烁定时器ID */
let cursorTimerId: number | null = null;

// ===================== 工具函数 =====================
/**
 * 从HTML内容中提取纯文本
 * @param {string} html - HTML字符串
 * @returns {string} 纯文本
 */
function extractPlainText(html: string): string {
  if (!browser) return html;
  
  const div = document.createElement('div');
  div.innerHTML = html;
  return div.textContent || div.innerText || '';
}

/**
 * 重建HTML，保持标签结构的同时截断到指定长度
 * @param {string} html - 原始HTML
 * @param {number} length - 要显示的字符长度
 * @returns {string} 截断后的HTML
 */
function rebuildHTML(html: string, length: number): string {
  if (!browser) return html.substring(0, length);
  
  const div = document.createElement('div');
  div.innerHTML = html;
  
  let currentLength = 0;
  const result = document.createElement('div');
  
  function processNode(node: Node, container: Element) {
    if (currentLength >= length) return false;
    
    if (node.nodeType === Node.TEXT_NODE) {
      const text = node.textContent || '';
      const remainingLength = length - currentLength;
      
      if (text.length <= remainingLength) {
        container.appendChild(node.cloneNode(true));
        currentLength += text.length;
      } else {
        const truncatedText = text.substring(0, remainingLength);
        const textNode = document.createTextNode(truncatedText);
        container.appendChild(textNode);
        currentLength += truncatedText.length;
        return false;
      }
    } else if (node.nodeType === Node.ELEMENT_NODE) {
      const element = node.cloneNode(false);
      container.appendChild(element);
      
      for (let child of node.childNodes) {
        if (!processNode(child, element)) {
          break;
        }
      }
    }
    
    return true;
  }
  
  for (let child of div.childNodes) {
    if (!processNode(child as Element, result)) {
      break;
    }
  }
  
  return result.innerHTML;
}

/**
 * 开始打字机动画
 * @returns {void}
 */
function startAnimation(): void {
  if (isPlaying || !plainText) return;
  
  isPlaying = true;
  isComplete = false;
  currentIndex = 0;
  displayedText = '';
  
  timerId = setInterval(() => {
    if (currentIndex < plainText.length) {
      currentIndex++;
      displayedText = rebuildHTML(content, currentIndex);
    } else {
      isComplete = true;
      isPlaying = false;
      clearInterval(timerId);
      
      if (loop) {
        setTimeout(() => {
          reset();
          startAnimation();
        }, 1000);
      }
    }
  }, speed);
}

/**
 * 停止动画
 * @returns {void}
 */
function stopAnimation(): void {
  if (timerId !== null) {
    clearInterval(timerId);
    timerId = null;
  }
  isPlaying = false;
}

/**
 * 重置动画状态
 * @returns {void}
 */
function reset(): void {
  stopAnimation();
  currentIndex = 0;
  displayedText = '';
  isComplete = false;
}

/**
 * 完成动画（立即显示全部内容）
 * @returns {void}
 */
function complete(): void {
  stopAnimation();
  currentIndex = plainText.length;
  displayedText = content;
  isComplete = true;
}

/**
 * 开始光标闪烁
 * @returns {void}
 */
function startCursorBlink(): void {
  if (!showCursor) return;
  
  cursorTimerId = setInterval(() => {
    cursorVisible = !cursorVisible;
  }, 500);
}

/**
 * 停止光标闪烁
 * @returns {void}
 */
function stopCursorBlink(): void {
  if (cursorTimerId !== null) {
    clearInterval(cursorTimerId);
    cursorTimerId = null;
  }
  cursorVisible = true;
}

// ===================== 响应式更新 =====================
$: {
  if (content) {
    plainText = extractPlainText(content);
    if (autoStart) {
      reset();
      startAnimation();
    }
  }
}

// ===================== 生命周期 =====================
onMount(() => {
  if (showCursor) {
    startCursorBlink();
  }
  
  if (autoStart && content) {
    plainText = extractPlainText(content);
    startAnimation();
  }
});

onDestroy(() => {
  stopAnimation();
  stopCursorBlink();
});

// ===================== 导出控制方法 =====================
export { startAnimation as start, stopAnimation as stop, reset, complete };
</script>

<!-- ===================== 组件模板 ===================== -->
<div class="typewriter-container relative inline-block">
  <!-- 文本内容 -->
  <span class="typewriter-text">
    {@html displayedText}
  </span>
  
  <!-- 光标 -->
  {#if showCursor && (!isComplete || isPlaying)}
    <span 
      class="typewriter-cursor inline-block w-0.5 h-5 bg-current ml-0.5 transition-opacity duration-100 {cursorVisible ? 'opacity-100' : 'opacity-0'}"
      aria-hidden="true"
    ></span>
  {/if}
  
  <!-- 控制按钮（可选） -->
  <div class="typewriter-controls mt-2 flex gap-2 text-xs opacity-0 group-hover:opacity-100 transition-opacity duration-300">
    {#if !autoStart}
      <button 
        class="px-2 py-1 bg-blue-500 text-white rounded hover:bg-blue-600 transition-colors"
        on:click={startAnimation}
        disabled={isPlaying}
      >
        {isPlaying ? '播放中...' : '开始'}
      </button>
    {/if}
    
    {#if isPlaying}
      <button 
        class="px-2 py-1 bg-red-500 text-white rounded hover:bg-red-600 transition-colors"
        on:click={stopAnimation}
      >
        停止
      </button>
    {/if}
    
    {#if !isComplete && currentIndex > 0}
      <button 
        class="px-2 py-1 bg-green-500 text-white rounded hover:bg-green-600 transition-colors"
        on:click={complete}
      >
        完成
      </button>
    {/if}
    
    {#if isComplete}
      <button 
        class="px-2 py-1 bg-gray-500 text-white rounded hover:bg-gray-600 transition-colors"
        on:click={() => { reset(); startAnimation(); }}
      >
        重播
      </button>
    {/if}
  </div>
</div>

<style>
  .typewriter-container {
    /* 确保容器有足够的空间 */
    min-height: 1.25rem;
  }
  
  .typewriter-cursor {
    /* 光标动画 */
    animation: blink 1s infinite;
  }
  
  @keyframes blink {
    0%, 50% { opacity: 1; }
    51%, 100% { opacity: 0; }
  }
  
  /* 悬停时显示控制按钮 */
  .typewriter-container:hover .typewriter-controls {
    opacity: 1;
  }
</style>

<!--
  关键算法说明：
  - extractPlainText: 从HTML中提取纯文本用于计算动画长度
  - rebuildHTML: 重建HTML结构，只显示指定长度的文本内容
  - 支持HTML高亮内容的逐字显示
  待优化事项：
  - 支持更多动画效果（淡入、滑动等）
  - 支持音效
  - 支持暂停/恢复功能
  兼容性说明：
  - 现代浏览器支持，使用DOM API处理HTML内容
-->