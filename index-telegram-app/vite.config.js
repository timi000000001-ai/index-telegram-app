import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [sveltekit()],
	build: {
		// 压缩配置
		minify: 'terser',
		// Terser 压缩选项
		terserOptions: {
			compress: {
				drop_console: true, // 移除console语句
				drop_debugger: true, // 移除debugger语句
				pure_funcs: ['console.log'] // 移除指定函数调用
			}
		},
		// 代码分割配置
		rollupOptions: {
			output: {
				// 最小化文件分割，尽可能合并所有代码
				manualChunks: undefined, // 禁用手动代码分割
				// 简化文件名
				chunkFileNames: '[name]-[hash].js',
				entryFileNames: '[name]-[hash].js',
				assetFileNames: '[name]-[hash].[ext]'
			},
			// Tree shaking 优化
			treeshake: {
				preset: 'recommended'
			}
		},
		// 设置chunk大小警告阈值
		chunkSizeWarningLimit: 1000,
		// 资源内联阈值
		assetsInlineLimit: 4096,
		// 启用CSS代码分割
		cssCodeSplit: false
	},
	// 优化依赖预构建
	optimizeDeps: {
		include: ['svelte'],
		// 强制预构建依赖
		force: false
	},
	// 预览配置
	preview: {
		port: 4173
	},
	// 开发服务器配置
	server: {
		fs: {
			// 允许为项目根目录以上的文件提供服务
			strict: false
		},

	}
});
