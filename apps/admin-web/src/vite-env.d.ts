/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_URL: string;
  readonly VITE_GRAPHQL_URL: string;
  readonly VITE_WS_URL: string;
  readonly VITE_HASURA_ADMIN_SECRET: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
