import { createRootRoute, Outlet } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";

export const Route = createRootRoute({
	component: () => (
		<div>
			<Outlet />
			<TanStackRouterDevtools />
		</div >
	),
	notFoundComponent: () => <div>404 Not Found</div>,
})
