import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { createConnectTransport } from "@connectrpc/connect-web";
import { TransportProvider } from "@connectrpc/connect-query";
import { createRouter, RouterProvider } from '@tanstack/react-router'
import { routeTree } from './routeTree.gen'
import { ReactQueryDevtools } from '@tanstack/react-query-devtools'
import { Toaster } from 'react-hot-toast'
import { GetBaseUrl } from './lib/baseURL';
import { GetJwtToken } from './lib/auth';

import gsap from 'gsap';
import { useGSAP } from '@gsap/react';

gsap.registerPlugin(useGSAP);

const router = createRouter({ routeTree })


const finalTransport = createConnectTransport({
	baseUrl: GetBaseUrl(),
	interceptors: [
		(next) => (request) => {
			const token = GetJwtToken();
			request.header.append("Authorization", `Bearer ${token}`)
			return next(request)
		},
	],
});

const queryClient = new QueryClient()

// Register the router instance for type safety
declare module '@tanstack/react-router' {
	interface Register {
		router: typeof router
	}
}
const rootElement = document.getElementById('root')!
if (!rootElement.innerHTML) {
	const root = createRoot(rootElement)
	root.render(
		<StrictMode>
			<Toaster toastOptions={{ duration: 1000 }} />
			<TransportProvider transport={finalTransport}>
				<QueryClientProvider client={queryClient}>
					<ReactQueryDevtools buttonPosition='top-left' />
					<RouterProvider router={router} />
				</QueryClientProvider>
			</TransportProvider >
		</StrictMode>,
	)
}

