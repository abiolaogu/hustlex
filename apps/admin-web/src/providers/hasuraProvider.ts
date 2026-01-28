import dataProviderHasura from "@refinedev/hasura";
import { GraphQLClient } from "graphql-request";

// Get config from runtime or build-time environment
const getConfig = () => {
  // Runtime config (from Docker)
  const runtimeConfig = (window as any).__ENV__ || {};

  return {
    apiUrl: runtimeConfig.VITE_GRAPHQL_URL ||
            import.meta.env?.VITE_GRAPHQL_URL ||
            "http://localhost:8080/v1/graphql",
    adminSecret: runtimeConfig.VITE_HASURA_ADMIN_SECRET ||
                 import.meta.env?.VITE_HASURA_ADMIN_SECRET ||
                 "hustlex_hasura_admin_secret",
  };
};

const config = getConfig();

const client = new GraphQLClient(config.apiUrl, {
  headers: {
    "x-hasura-admin-secret": config.adminSecret,
  },
});

// Cast to any to avoid version mismatch issues
export const dataProvider = dataProviderHasura(client as any);

export const liveProvider = {
  subscribe: () => {
    // Placeholder for live subscription
    return () => {};
  },
  unsubscribe: () => {},
};
