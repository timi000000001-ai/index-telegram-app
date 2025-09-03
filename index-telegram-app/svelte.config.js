import adapter from '@sveltejs/adapter-static';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	kit: {
		// 使用静态适配器生成HTML页面
		adapter: adapter({
			pages: 'build',
			assets: 'build',
			fallback: 'index.html',
			precompress: true,
			// 启用严格模式
			strict: true
		}),
		// 启用预渲染以减少运行时文件
		prerender: {
			entries: ['*'],
			// 启用并发预渲染
			concurrency: 4
		},
		// 服务工作者配置
		serviceWorker: {
			register: false
		},
		// 内联样式阈值
		inlineStyleThreshold: 1024,
		// CSP配置
		csp: {
			mode: 'auto'
		}
	},
	// 编译器选项
	compilerOptions: {
		// 启用hydratable
		hydratable: true,
		// 启用CSS变量
		cssHash: ({ hash, css }) => `svelte-${hash(css)}`
	}
};

export default config;
