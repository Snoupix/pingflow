/// <reference types="vite/client" />

interface ImportMetaEnv {
    readonly VITE_WEBSOCKET_URI: string;
    readonly VITE_WEBSOCKET_PORT: number;
    readonly VITE_WEBSOCKET_TIMEOUT: number;
    readonly VITE_WEBSOCKET_RECONNECTION_TIMES: number;
}

interface ImportMeta {
    readonly env: ImportMetaEnv;
}
