<script>
	import { onMount } from 'svelte';
	import { apiFetch } from '$lib/api.js';

	// 类型定义
	/** @typedef {Object} Bot
	 * @property {number} id
	 * @property {string} name
	 * @property {string} status
	 * @property {string} uptime
	 * @property {number} messages
	 * @property {string} responseTime
	 * @property {number} errors
	 * @property {string} lastActivity
	 * @property {string} createdAt
	 */

	// 机器人状态数据
	/** @type {Bot[]} */
	let bots = [];
	let loading = true;
	/** @type {string|null} */
	let error = null;
	let activeTab = 'status'; // status, clone, config

	// 统计数据
	let stats = {
		online: 0,
		total: 0,
		totalMessages: 0,
		totalUptime: '0天'
	};

	// 性能数据（模拟）
	let performanceData = {
		responseTime: [120, 95, 110, 85, 105, 90, 115, 100],
		messageCount: [150, 120, 180, 140, 200, 100, 160, 130]
	};

	// 错误日志
	let errorLogs = [
		{ time: '10:30', type: 'warning', message: '连接超时警告' },
		{ time: '09:15', type: 'error', message: 'API调用失败' },
		{ time: '08:45', type: 'info', message: '重连成功' }
	];

	// 获取机器人列表
	async function fetchBots() {
		try {
			loading = true;
			const response = await apiFetch('/api/bots/status');
			
			if (response.success) {
				bots = response.bots || [
					{
						id: 1,
						name: '搜索机器人 #1',
						status: 'online',
						uptime: '2天5小时',
						messages: 1500,
						responseTime: '120ms',
						errors: 2,
						lastActivity: '2分钟前',
						createdAt: '2025-01-20'
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
						createdAt: '2025-01-19'
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
						createdAt: '2025-01-22'
					}
				];
				updateStats();
			} else {
				error = response.error || '获取机器人列表失败';
			}
		} catch (err) {
			error = '网络错误，请稍后重试';
			console.error('获取机器人列表失败:', err);
		} finally {
			loading = false;
		}
	}

	// 更新统计数据
	function updateStats() {
		stats.total = bots.length;
		stats.online = bots.filter(bot => bot.status === 'online').length;
		stats.totalMessages = bots.reduce((sum, bot) => sum + bot.messages, 0);
	}

	// 机器人操作
	/**
	 * @param {number} botId
	 */
	async function cloneBot(botId) {
		alert(`克隆机器人 ID: ${botId}`);
	}

	/**
	 * @param {number} botId
	 */
	async function configBot(botId) {
		alert(`配置机器人 ID: ${botId}`);
	}

	/**
	 * @param {number} botId
	 */
	async function stopBot(botId) {
		if (confirm('确定要停止这个机器人吗？')) {
			alert(`停止机器人 ID: ${botId}`);
		}
	}



	// 获取状态颜色
	/**
	 * @param {string} status
	 */
	function getStatusColor(status) {
		switch (status) {
			case 'online': return 'text-green-600';
			case 'starting': return 'text-yellow-600';
			case 'offline': return 'text-red-600';
			default: return 'text-gray-600';
		}
	}

	/**
	 * @param {string} status
	 */
	function getStatusBgColor(status) {
		switch (status) {
			case 'online': return 'bg-green-100';
			case 'starting': return 'bg-yellow-100';
			case 'offline': return 'bg-red-100';
			default: return 'bg-gray-100';
		}
	}

	/**
	 * @param {string} status
	 */
	function getStatusDotColor(status) {
		switch (status) {
			case 'online': return 'bg-green-500';
			case 'starting': return 'bg-yellow-500';
			case 'offline': return 'bg-red-500';
			default: return 'bg-gray-500';
		}
	}

	onMount(() => {
		fetchBots();
	});
</script>

<svelte:head>
	<title>机器人管理 - Telegram Search</title>
</svelte:head>

<div class="min-h-screen bg-gradient-to-br from-slate-50 to-slate-100">
	<!-- 机器人概览区域 -->
	<div class="bg-white/80 backdrop-blur-sm border-b border-slate-200/60 p-6 shadow-sm">
		<div class="max-w-7xl mx-auto">
			<h1 class="text-3xl font-bold bg-gradient-to-r from-slate-800 to-slate-600 bg-clip-text text-transparent mb-6">机器人概览</h1>
			
			<div class="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-4 gap-6 mb-8">
				<!-- 状态卡片 -->
				{#each bots.slice(0, 3) as bot}
					<div class="bg-white/90 backdrop-blur-sm border border-slate-200/60 rounded-xl p-5 shadow-lg hover:shadow-xl transition-all duration-300 hover:-translate-y-1">
						<div class="flex items-center mb-3">
							<div class="w-3 h-3 rounded-full {getStatusDotColor(bot.status)} mr-3 shadow-sm"></div>
							<h3 class="font-semibold text-slate-800 text-sm">{bot.name}</h3>
						</div>
						<p class="text-xs text-slate-600 mb-2">运行时间: {bot.uptime}</p>
						<p class="text-xs text-slate-600">处理消息: {bot.messages}条</p>
					</div>
				{/each}
				
				<!-- 快捷操作按钮 -->
				<div class="flex flex-col gap-4">
					<button class="bg-gradient-to-r from-sky-500 to-sky-600 hover:from-sky-600 hover:to-sky-700 text-white px-5 py-3 rounded-xl text-sm font-medium transition-all duration-300 shadow-lg hover:shadow-xl hover:-translate-y-1">
						新建机器人
					</button>
					<button class="bg-gradient-to-r from-emerald-500 to-emerald-600 hover:from-emerald-600 hover:to-emerald-700 text-white px-5 py-3 rounded-xl text-sm font-medium transition-all duration-300 shadow-lg hover:shadow-xl hover:-translate-y-1">
						克隆机器人
					</button>
				</div>
			</div>
		</div>
	</div>

	<!-- 功能选项卡 -->
	<div class="bg-white/80 backdrop-blur-sm border-b border-slate-200/60 shadow-sm">
		<div class="max-w-7xl mx-auto px-6">
			<nav class="flex space-x-8">
				<button 
					class="py-4 px-2 border-b-2 font-medium text-sm transition-all duration-300
						{activeTab === 'status' ? 'border-sky-500 text-sky-600 bg-sky-50/50' : 'border-transparent text-slate-500 hover:text-slate-700 hover:border-slate-300 hover:bg-slate-50/50'}"
					on:click={() => activeTab = 'status'}
				>
					状态监控
				</button>
				<button 
					class="py-4 px-2 border-b-2 font-medium text-sm transition-all duration-300
						{activeTab === 'clone' ? 'border-sky-500 text-sky-600 bg-sky-50/50' : 'border-transparent text-slate-500 hover:text-slate-700 hover:border-slate-300 hover:bg-slate-50/50'}"
					on:click={() => activeTab = 'clone'}
				>
					克隆管理
				</button>
				<button 
					class="py-4 px-2 border-b-2 font-medium text-sm transition-all duration-300
						{activeTab === 'config' ? 'border-sky-500 text-sky-600 bg-sky-50/50' : 'border-transparent text-slate-500 hover:text-slate-700 hover:border-slate-300 hover:bg-slate-50/50'}"
					on:click={() => activeTab = 'config'}
				>
					配置设置
				</button>
			</nav>
		</div>
	</div>

	<!-- 内容区域 -->
	<div class="max-w-7xl mx-auto p-6">
		{#if activeTab === 'status'}
			<!-- 状态监控内容 -->
			<div class="space-y-8">
				<h2 class="text-2xl font-bold bg-gradient-to-r from-slate-800 to-slate-600 bg-clip-text text-transparent">实时状态监控</h2>
				
				<!-- 性能图表区域 -->
				<div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
					<!-- 响应时间趋势图 -->
					<div class="bg-white/90 backdrop-blur-sm border border-slate-200/60 rounded-xl p-6 shadow-lg hover:shadow-xl transition-all duration-300">
						<h3 class="text-lg font-semibold text-slate-800 mb-4 text-center">响应时间趋势图</h3>
						<div class="h-32 bg-gradient-to-br from-slate-50 to-slate-100 rounded-lg flex items-end justify-center space-x-2 p-4">
							{#each performanceData.responseTime as time, i}
								<div 
									class="bg-gradient-to-t from-sky-500 to-sky-400 rounded-t shadow-sm" 
									style="height: {(time / 150) * 100}%; width: 12px;"
									title="{time}ms"
								></div>
							{/each}
						</div>
					</div>

					<!-- 消息处理统计 -->
					<div class="bg-white/90 backdrop-blur-sm border border-slate-200/60 rounded-xl p-6 shadow-lg hover:shadow-xl transition-all duration-300">
						<h3 class="text-lg font-semibold text-slate-800 mb-4 text-center">消息处理统计</h3>
						<div class="h-32 bg-gradient-to-br from-slate-50 to-slate-100 rounded-lg flex items-end justify-center space-x-2 p-4">
							{#each performanceData.messageCount as count, i}
								<div 
									class="bg-gradient-to-t from-emerald-500 to-emerald-400 rounded-t shadow-sm" 
									style="height: {(count / 200) * 100}%; width: 12px;"
									title="{count}条"
								></div>
							{/each}
						</div>
					</div>

					<!-- 错误日志区域 -->
					<div class="bg-white/90 backdrop-blur-sm border border-slate-200/60 rounded-xl p-6 shadow-lg hover:shadow-xl transition-all duration-300">
						<h3 class="text-lg font-semibold text-slate-800 mb-4 text-center">最近错误日志</h3>
						<div class="space-y-3">
							{#each errorLogs as log}
								<div class="p-3 rounded-lg border transition-all duration-200 hover:shadow-md
									{log.type === 'warning' ? 'bg-amber-50/80 border-amber-200/60 hover:bg-amber-50' : 
									 log.type === 'error' ? 'bg-red-50/80 border-red-200/60 hover:bg-red-50' : 
									 'bg-sky-50/80 border-sky-200/60 hover:bg-sky-50'}
								">
									<p class="text-xs font-medium
										{log.type === 'warning' ? 'text-amber-800' : 
										 log.type === 'error' ? 'text-red-800' : 
										 'text-sky-800'}
									">
										{log.time} - {log.message}
									</p>
								</div>
							{/each}
						</div>
					</div>
				</div>

				<!-- 详细状态表格 -->
				<div class="bg-white/90 backdrop-blur-sm border border-slate-200/60 rounded-xl overflow-hidden shadow-lg">
					<div class="px-6 py-5 border-b border-slate-200/60 bg-gradient-to-r from-slate-50/50 to-white/50">
						<h3 class="text-lg font-semibold text-slate-800">机器人详细状态</h3>
					</div>
					
					{#if loading}
						<div class="p-12 text-center">
							<div class="animate-spin rounded-full h-10 w-10 border-b-2 border-sky-500 mx-auto"></div>
							<p class="mt-3 text-slate-600 font-medium">加载中...</p>
						</div>
					{:else if error}
						<div class="p-12 text-center">
							<p class="text-red-600 font-medium">{error}</p>
							<button 
								class="mt-4 bg-gradient-to-r from-sky-500 to-sky-600 hover:from-sky-600 hover:to-sky-700 text-white px-5 py-2.5 rounded-lg font-medium transition-all duration-300 shadow-lg hover:shadow-xl"
								on:click={fetchBots}
							>
								重试
							</button>
						</div>
					{:else}
						<div class="overflow-x-auto">
							<table class="min-w-full divide-y divide-slate-200/60">
								<thead class="bg-gradient-to-r from-slate-50/80 to-slate-100/80">
									<tr>
										<th class="px-6 py-4 text-left text-xs font-semibold text-slate-600 uppercase tracking-wider">机器人名称</th>
										<th class="px-6 py-4 text-left text-xs font-semibold text-slate-600 uppercase tracking-wider">状态</th>
										<th class="px-6 py-4 text-left text-xs font-semibold text-slate-600 uppercase tracking-wider">运行时间</th>
										<th class="px-6 py-4 text-left text-xs font-semibold text-slate-600 uppercase tracking-wider">处理消息</th>
										<th class="px-6 py-4 text-left text-xs font-semibold text-slate-600 uppercase tracking-wider">响应时间</th>
										<th class="px-6 py-4 text-left text-xs font-semibold text-slate-600 uppercase tracking-wider">错误数</th>
										<th class="px-6 py-4 text-left text-xs font-semibold text-slate-600 uppercase tracking-wider">最后活动</th>
										<th class="px-6 py-4 text-left text-xs font-semibold text-slate-600 uppercase tracking-wider">操作</th>
									</tr>
								</thead>
								<tbody class="bg-white/50 divide-y divide-slate-200/40">
									{#each bots as bot, i}
										<tr class="{i % 2 === 0 ? 'bg-white/30' : 'bg-slate-50/30'} hover:bg-slate-50/60 transition-colors duration-200">
											<td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-slate-800">{bot.name}</td>
											<td class="px-6 py-4 whitespace-nowrap">
												<div class="flex items-center">
													<div class="w-2.5 h-2.5 rounded-full {getStatusDotColor(bot.status)} mr-3 shadow-sm"></div>
													<span class="text-sm font-medium {getStatusColor(bot.status)}">
														{bot.status === 'online' ? '在线' : bot.status === 'starting' ? '启动中' : '离线'}
													</span>
												</div>
											</td>
											<td class="px-6 py-4 whitespace-nowrap text-sm text-slate-700">{bot.uptime}</td>
											<td class="px-6 py-4 whitespace-nowrap text-sm text-slate-700">{bot.messages}</td>
											<td class="px-6 py-4 whitespace-nowrap text-sm text-slate-700">{bot.responseTime}</td>
											<td class="px-6 py-4 whitespace-nowrap text-sm text-slate-700">{bot.errors}</td>
											<td class="px-6 py-4 whitespace-nowrap text-sm text-slate-700">{bot.lastActivity}</td>
											<td class="px-6 py-4 whitespace-nowrap text-sm font-medium space-x-2">
												<button 
													class="bg-gradient-to-r from-emerald-500 to-emerald-600 hover:from-emerald-600 hover:to-emerald-700 text-white px-3 py-1.5 rounded-lg text-xs font-medium transition-all duration-300 shadow-sm hover:shadow-md"
													on:click={() => cloneBot(bot.id)}
												>
													克隆
												</button>
												<button 
													class="bg-gradient-to-r from-sky-500 to-sky-600 hover:from-sky-600 hover:to-sky-700 text-white px-3 py-1.5 rounded-lg text-xs font-medium transition-all duration-300 shadow-sm hover:shadow-md"
													on:click={() => configBot(bot.id)}
												>
													配置
												</button>
												<button 
													class="bg-gradient-to-r from-red-500 to-red-600 hover:from-red-600 hover:to-red-700 text-white px-3 py-1.5 rounded-lg text-xs font-medium transition-all duration-300 shadow-sm hover:shadow-md"
													on:click={() => stopBot(bot.id)}
												>
													停止
												</button>
											</td>
										</tr>
									{/each}
								</tbody>
							</table>
						</div>
					{/if}
				</div>
			</div>
		{:else if activeTab === 'clone'}
			<!-- 克隆管理内容 -->
			<div class="space-y-8">
				<div>
					<h2 class="text-2xl font-bold bg-gradient-to-r from-slate-800 to-slate-600 bg-clip-text text-transparent mb-2">克隆机器人</h2>
					<p class="text-slate-600">复制现有机器人配置，快速创建新的机器人实例</p>
				</div>
				
				<!-- 步骤指示器 -->
				<div class="bg-white/90 backdrop-blur-sm border border-slate-200/60 rounded-xl p-6 shadow-lg">
					<div class="flex items-center justify-center space-x-8">
						<div class="flex items-center">
							<div class="w-10 h-10 bg-gradient-to-r from-sky-500 to-sky-600 text-white rounded-full flex items-center justify-center text-sm font-bold shadow-lg">1</div>
							<span class="ml-3 text-sm font-semibold text-sky-600">选择源机器人</span>
						</div>
						<div class="w-16 h-1 bg-gradient-to-r from-sky-500 to-sky-600 rounded-full"></div>
						<div class="flex items-center">
							<div class="w-10 h-10 bg-gradient-to-r from-slate-400 to-slate-500 text-white rounded-full flex items-center justify-center text-sm font-bold shadow-lg">2</div>
							<span class="ml-3 text-sm font-semibold text-slate-500">配置克隆选项</span>
						</div>
					</div>
				</div>
				
				<div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
					<!-- 左侧：源机器人选择 -->
					<div class="bg-white/90 backdrop-blur-sm border border-slate-200/60 rounded-xl p-6 shadow-lg">
						<h3 class="text-lg font-semibold text-slate-800 mb-5">选择要克隆的机器人</h3>
						
						<!-- 搜索框 -->
						<div class="mb-6">
							<div class="relative">
								<input 
									type="text" 
									placeholder="搜索机器人名称..."
									class="w-full pl-4 pr-12 py-3 border border-slate-300/60 rounded-xl focus:ring-2 focus:ring-sky-500/50 focus:border-sky-500 bg-white/80 backdrop-blur-sm transition-all duration-300"
								/>
								<div class="absolute inset-y-0 right-0 pr-4 flex items-center">
									<svg class="h-5 w-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
									</svg>
								</div>
							</div>
						</div>
						
						<h4 class="font-semibold text-slate-800 mb-4">可用机器人列表</h4>
						<div class="space-y-4">
							{#each bots as bot}
								<div class="border-2 border-sky-500/60 bg-gradient-to-r from-sky-50/80 to-sky-100/60 rounded-xl p-5 cursor-pointer hover:shadow-lg transition-all duration-300 hover:-translate-y-1">
									<div class="flex items-center justify-between">
										<div class="flex items-center">
											<div class="w-3 h-3 rounded-full {getStatusDotColor(bot.status)} mr-4 shadow-sm"></div>
											<div>
												<h5 class="font-semibold text-slate-800 mb-1">{bot.name}</h5>
												<p class="text-xs text-slate-600 mb-1">创建时间: {bot.createdAt}</p>
												<p class="text-xs text-slate-600 mb-1">处理消息: {bot.messages}条 | 运行时间: {bot.uptime}</p>
												<p class="text-xs font-medium {getStatusColor(bot.status)}">状态: {bot.status === 'online' ? '运行中' : bot.status === 'starting' ? '启动中' : '离线'}</p>
											</div>
										</div>
										<div class="w-7 h-7 bg-gradient-to-r from-sky-500 to-sky-600 text-white rounded-full flex items-center justify-center shadow-lg">
											<svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
												<path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"></path>
											</svg>
										</div>
									</div>
								</div>
							{/each}
						</div>
					</div>
					
					<!-- 右侧：克隆配置选项 -->
					<div class="bg-white/90 backdrop-blur-sm border border-slate-200/60 rounded-xl p-6 shadow-lg">
						<h3 class="text-lg font-semibold text-slate-800 mb-5">克隆配置选项</h3>
						
						<!-- 基本信息 -->
						<div class="mb-8">
							<h4 class="font-semibold text-slate-800 mb-4">基本信息</h4>
							<div class="space-y-5">
								<div>
								<label for="bot-name" class="block text-sm font-medium text-slate-700 mb-2">新机器人名称</label>
								<input 
									id="bot-name"
									type="text" 
									value="搜索机器人 #1 - 副本"
									class="w-full px-4 py-3 border border-slate-300/60 rounded-xl focus:ring-2 focus:ring-sky-500/50 focus:border-sky-500 bg-white/80 backdrop-blur-sm transition-all duration-300"
								/>
								</div>
								<div>
								<label for="bot-description" class="block text-sm font-medium text-slate-700 mb-2">描述信息</label>
								<textarea 
									id="bot-description"
									rows="3"
									class="w-full px-4 py-3 border border-slate-300/60 rounded-xl focus:ring-2 focus:ring-sky-500/50 focus:border-sky-500 bg-white/80 backdrop-blur-sm transition-all duration-300 resize-none"
									placeholder="基于搜索机器人 #1 克隆的新实例"
								></textarea>
								</div>
							</div>
						</div>
						
						<!-- 克隆选项 -->
						<div class="mb-8">
							<h4 class="font-semibold text-slate-800 mb-4">克隆选项</h4>
							<div class="space-y-4">
								<label class="flex items-center p-3 rounded-lg hover:bg-slate-50/50 transition-colors duration-200 cursor-pointer">
									<input type="checkbox" checked class="rounded border-slate-300 text-sky-600 focus:ring-sky-500/50 w-4 h-4">
									<span class="ml-3 text-sm font-medium text-slate-800">搜索配置 (搜索算法、匹配规则)</span>
								</label>
								<label class="flex items-center p-3 rounded-lg hover:bg-slate-50/50 transition-colors duration-200 cursor-pointer">
									<input type="checkbox" checked class="rounded border-slate-300 text-sky-600 focus:ring-sky-500/50 w-4 h-4">
									<span class="ml-3 text-sm font-medium text-slate-800">过滤规则 (时间、来源、类型过滤器)</span>
								</label>
								<label class="flex items-center p-3 rounded-lg hover:bg-slate-50/50 transition-colors duration-200 cursor-pointer">
									<input type="checkbox" checked class="rounded border-slate-300 text-sky-600 focus:ring-sky-500/50 w-4 h-4">
									<span class="ml-3 text-sm font-medium text-slate-800">响应模板 (消息格式、错误提示)</span>
								</label>
								<label class="flex items-center p-3 rounded-lg hover:bg-slate-50/50 transition-colors duration-200 cursor-pointer">
									<input type="checkbox" class="rounded border-slate-300 text-sky-600 focus:ring-sky-500/50 w-4 h-4">
									<span class="ml-3 text-sm font-medium text-slate-800">权限设置 (访问控制、管理权限)</span>
								</label>
								<label class="flex items-center p-3 rounded-lg hover:bg-slate-50/50 transition-colors duration-200 cursor-pointer">
									<input type="checkbox" class="rounded border-slate-300 text-sky-600 focus:ring-sky-500/50 w-4 h-4">
									<span class="ml-3 text-sm font-medium text-slate-800">数据源连接 (群组、频道连接)</span>
								</label>
								<label class="flex items-center p-3 rounded-lg hover:bg-slate-50/50 transition-colors duration-200 cursor-pointer">
									<input type="checkbox" checked class="rounded border-slate-300 text-sky-600 focus:ring-sky-500/50 w-4 h-4">
									<span class="ml-3 text-sm font-medium text-slate-800">性能配置 (CPU、内存、存储设置)</span>
								</label>
							</div>
						</div>
						
						<!-- 操作按钮 -->
						<div class="flex justify-end space-x-4 pt-6 border-t border-slate-200/60">
							<button class="px-6 py-3 text-slate-600 border border-slate-300/60 rounded-xl hover:bg-slate-50/50 font-medium transition-all duration-300 shadow-sm">
								取消
							</button>
							<button class="px-6 py-3 bg-gradient-to-r from-sky-500 to-blue-600 text-white rounded-xl hover:from-sky-600 hover:to-blue-700 font-medium transition-all duration-300 shadow-lg hover:shadow-xl">
								开始克隆
							</button>
						</div>

					</div>
				</div>
			</div>
		{:else if activeTab === 'config'}
			<!-- 配置设置内容 -->
			<div class="space-y-8">
				<h2 class="text-xl font-bold text-slate-800">配置设置</h2>
				<div class="bg-white/90 backdrop-blur-sm border border-slate-200/60 rounded-xl p-8 shadow-lg">
					<div class="text-center py-12">
						<div class="w-16 h-16 mx-auto mb-4 bg-gradient-to-br from-sky-100 to-blue-100 rounded-full flex items-center justify-center">
							<svg class="w-8 h-8 text-sky-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"></path>
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
							</svg>
						</div>
						<h3 class="text-lg font-semibold text-slate-800 mb-2">配置设置</h3>
						<p class="text-slate-600 mb-6">高级配置功能正在开发中，敬请期待...</p>
						<div class="inline-flex items-center px-4 py-2 bg-gradient-to-r from-sky-50 to-blue-50 border border-sky-200/60 rounded-lg text-sm text-sky-700">
							<svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
							</svg>
							即将推出更多配置选项
						</div>
					</div>
				</div>
			</div>
		{/if}
	</div>


</div>