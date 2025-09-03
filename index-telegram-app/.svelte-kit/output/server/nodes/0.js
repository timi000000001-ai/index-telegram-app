import * as universal from '../entries/pages/_layout.js';

export const index = 0;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/_layout.svelte.js')).default;
export { universal };
export const universal_id = "src/routes/+layout.js";
export const imports = ["_app/immutable/nodes/0.CVUjmk3Y.js","_app/immutable/chunks/Cm9PB-7n.js","_app/immutable/chunks/ffK739sN.js","_app/immutable/chunks/BnjqhL1k.js","_app/immutable/chunks/SHOpF_Wr.js","_app/immutable/chunks/BlKWISPa.js","_app/immutable/chunks/BYfX5H5R.js","_app/immutable/chunks/DvxUEipc.js","_app/immutable/chunks/f_3be0LC.js"];
export const stylesheets = ["_app/immutable/assets/0.DQuk6lQC.css"];
export const fonts = [];
