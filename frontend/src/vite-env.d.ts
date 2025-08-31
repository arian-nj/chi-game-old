/// <reference types="vite/client" />

declare module '*.lottie' {
  const src: string;
  export default src;
}
declare module '*.wasm' {
  const url: string;
  export default url;
}
