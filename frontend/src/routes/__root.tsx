import { createRootRoute, Outlet } from "@tanstack/react-router";

export const Route = createRootRoute({
	component: () => (
		<div>
			<Outlet />
			{/* <TanStackRouterDevtools /> */}
		</div >
	),
	notFoundComponent: () => <div>404 Not Found</div>,
})
