import { createConnectTransport } from "@connectrpc/connect-web";
import { GetBaseUrl } from "./baseURL";
import { GetJwtToken } from "./auth";

export const rawTransport = createConnectTransport({
  baseUrl: GetBaseUrl(),
})
export const authTransport = createConnectTransport({
  baseUrl: GetBaseUrl(),
  interceptors: [
    (next) => (request) => {
      const token = GetJwtToken();
      request.header.append("Authorization", `Bearer ${token}`)
      return next(request)
    },
  ],
})
