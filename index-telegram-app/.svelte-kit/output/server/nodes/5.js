

export const index = 5;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/daily-stats/_page.svelte.js')).default;
export const imports = ["_app/immutable/nodes/5.MQV4TwdP.js","_app/immutable/chunks/Cm9PB-7n.js","_app/immutable/chunks/ffK739sN.js","_app/immutable/chunks/BnjqhL1k.js","_app/immutable/chunks/D7lVvpcI.js","_app/immutable/chunks/BxFA2AXj.js","_app/immutable/chunks/SHOpF_Wr.js","_app/immutable/chunks/BlKWISPa.js"];
export const stylesheets = ["_app/immutable/assets/5.PQP-262y.css"];
export const fonts = [];
