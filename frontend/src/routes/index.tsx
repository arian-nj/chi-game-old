import { createFileRoute } from "@tanstack/react-router";


export const Route = createFileRoute("/")({
	component: Index,
})

function Index() {
	return (
		<div className="p-2">
			<button
				type="button"
				className="bg-gradient-to-r from-pink-500 to-orange-400 rounded-lg text-center
				px-7 py-3 text-6xl text-gray-100 text-shadow-md
					focus:ring-4 focus:ring-pink-200 focus:outline-none focus:bg-gradient-to-tr
				"
			>Play</button>
		</div>
	)
}
