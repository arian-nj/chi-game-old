import { createConnectTransport } from "@connectrpc/connect-web";
import { GetApiUrl } from "./baseURL";
import { GetJwtToken } from "./auth";

export const rawTransport = createConnectTransport({
  baseUrl: GetApiUrl(),
})
export const authTransport = createConnectTransport({
  baseUrl: GetApiUrl(),
  interceptors: [
    (next) => (request) => {
      const token = GetJwtToken();
      request.header.append("Authorization", `Bearer ${token}`)
      return next(request)
    },
  ],
})
