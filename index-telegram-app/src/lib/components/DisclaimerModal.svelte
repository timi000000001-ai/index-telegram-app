<script>
  import { createEventDispatcher } from 'svelte';

  export let show = false;

  const dispatch = createEventDispatcher();

  function close() {
    show = false;
    dispatch('close');
  }

  function handleKeydown(/** @type {KeyboardEvent} */ event) {
    if (event.key === 'Escape') {
      close();
    }
  }
</script>

{#if show}
  <div
    class="fixed inset-0 z-50 flex items-center justify-center p-4"
    aria-modal="true"
    role="dialog"
    tabindex="-1"
    on:keydown={handleKeydown}
  >
    <div class="fixed inset-0 bg-black bg-opacity-50" on:click={close}></div>
    <div
      class="relative bg-white rounded-lg shadow-xl max-w-md w-full mx-auto p-6 space-y-4 transform transition-all sm:my-8 sm:w-full sm:max-w-lg"
    >
      <h2 class="text-2xl font-bold text-gray-900">法律声明</h2>

      <div class="prose prose-sm max-w-none text-gray-700">
        <h3 class="text-lg font-semibold text-gray-800">使用限制</h3>
        <p>本项目不适用于中国大陆。Telegram 在中国大陆受到政府的访问限制，本项目的数据收集和处理活动可能违反当地法律法规。</p>

        <h3 class="text-lg font-semibold text-gray-800">免责声明</h3>
        <p>本项目开发人员对因使用不当、违反当地法律或数据隐私问题而导致的任何后果概不负责。用户应自行评估法律风险，并在必要时咨询法律专业人士。</p>

        <h3 class="text-lg font-semibold text-gray-800">建议</h3>
        <p>如果您位于中国大陆，请不要下载、安装或运行本项目。请寻找符合当地法规的替代方案。</p>
      </div>

      <div class="flex justify-end">
        <button
          type="button"
          class="inline-flex justify-center rounded-md border border-transparent bg-blue-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
          on:click={close}
        >
          关闭
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
  /* 可以添加一些过渡效果 */
  .transform {
    transition: all 0.3s ease-out;
  }
</style>