export const manifest = (() => {
function __memo(fn) {
	let value;
	return () => value ??= (value = fn());
}

return {
	appDir: "_app",
	appPath: "_app",
	assets: new Set(["robots.txt"]),
	mimeTypes: {".txt":"text/plain"},
	_: {
		client: {start:"_app/immutable/entry/start.CsKQ93LY.js",app:"_app/immutable/entry/app.CQYSYtZr.js",imports:["_app/immutable/entry/start.CsKQ93LY.js","_app/immutable/chunks/Dmm1jFPB.js","_app/immutable/chunks/vk1Hrpyt.js","_app/immutable/chunks/0-kjiHYL.js","_app/immutable/entry/app.CQYSYtZr.js","_app/immutable/chunks/vk1Hrpyt.js","_app/immutable/chunks/Cm9PB-7n.js","_app/immutable/chunks/BGXLXVlW.js","_app/immutable/chunks/BIuud_EK.js","_app/immutable/chunks/DliHbt_6.js","_app/immutable/chunks/0-kjiHYL.js"],stylesheets:[],fonts:[],uses_env_dynamic_public:false},
		nodes: [
			__memo(() => import('./nodes/0.js')),
			__memo(() => import('./nodes/1.js')),
			__memo(() => import('./nodes/2.js')),
			__memo(() => import('./nodes/3.js')),
			__memo(() => import('./nodes/4.js')),
			__memo(() => import('./nodes/5.js')),
			__memo(() => import('./nodes/6.js'))
		],
		remotes: {
			
		},
		routes: [
			{
				id: "/",
				pattern: /^\/$/,
				params: [],
				page: { layouts: [0,], errors: [1,], leaf: 2 },
				endpoint: null
			},
			{
				id: "/api/bots/status",
				pattern: /^\/api\/bots\/status\/?$/,
				params: [],
				page: null,
				endpoint: __memo(() => import('./entries/endpoints/api/bots/status/_server.js'))
			},
			{
				id: "/api/bot/status",
				pattern: /^\/api\/bot\/status\/?$/,
				params: [],
				page: null,
				endpoint: __memo(() => import('./entries/endpoints/api/bot/status/_server.js'))
			},
			{
				id: "/api/collect",
				pattern: /^\/api\/collect\/?$/,
				params: [],
				page: null,
				endpoint: __memo(() => import('./entries/endpoints/api/collect/_server.js'))
			},
			{
				id: "/api/login",
				pattern: /^\/api\/login\/?$/,
				params: [],
				page: null,
				endpoint: __memo(() => import('./entries/endpoints/api/login/_server.js'))
			},
			{
				id: "/api/search",
				pattern: /^\/api\/search\/?$/,
				params: [],
				page: null,
				endpoint: __memo(() => import('./entries/endpoints/api/search/_server.js'))
			},
			{
				id: "/api/suggestions",
				pattern: /^\/api\/suggestions\/?$/,
				params: [],
				page: null,
				endpoint: __memo(() => import('./entries/endpoints/api/suggestions/_server.js'))
			},
			{
				id: "/api/trending",
				pattern: /^\/api\/trending\/?$/,
				params: [],
				page: null,
				endpoint: __memo(() => import('./entries/endpoints/api/trending/_server.js'))
			},
			{
				id: "/api/verify",
				pattern: /^\/api\/verify\/?$/,
				params: [],
				page: null,
				endpoint: __memo(() => import('./entries/endpoints/api/verify/_server.js'))
			},
			{
				id: "/bots",
				pattern: /^\/bots\/?$/,
				params: [],
				page: { layouts: [0,], errors: [1,], leaf: 3 },
				endpoint: null
			},
			{
				id: "/daily-new",
				pattern: /^\/daily-new\/?$/,
				params: [],
				page: { layouts: [0,], errors: [1,], leaf: 4 },
				endpoint: null
			},
			{
				id: "/daily-stats",
				pattern: /^\/daily-stats\/?$/,
				params: [],
				page: { layouts: [0,], errors: [1,], leaf: 5 },
				endpoint: null
			},
			{
				id: "/login",
				pattern: /^\/login\/?$/,
				params: [],
				page: { layouts: [0,], errors: [1,], leaf: 6 },
				endpoint: null
			}
		],
		prerendered_routes: new Set([]),
		matchers: async () => {
			
			return {  };
		},
		server_assets: {}
	}
}
})();
