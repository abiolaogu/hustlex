import dataProviderHasura from "@refinedev/hasura";
import { GraphQLClient } from "graphql-request";
import { createClient } from "graphql-ws";

const API_URL = import.meta.env.VITE_GRAPHQL_URL || "http://localhost:8080/v1/graphql";
const WS_URL = import.meta.env.VITE_WS_URL || "ws://localhost:8080/v1/graphql";

const client = new GraphQLClient(API_URL, {
  headers: {
    "x-hasura-admin-secret": import.meta.env.VITE_HASURA_ADMIN_SECRET || "hustlex_hasura_admin_secret",
  },
});

const wsClient = createClient({
  url: WS_URL,
  connectionParams: {
    headers: {
      "x-hasura-admin-secret": import.meta.env.VITE_HASURA_ADMIN_SECRET || "hustlex_hasura_admin_secret",
    },
  },
});

export const dataProvider = dataProviderHasura(client);

export const liveProvider = {
  subscribe: ({ channel, types, params, callback }: any) => {
    // Implement live subscription using wsClient
    return () => {};
  },
  unsubscribe: () => {},
};
