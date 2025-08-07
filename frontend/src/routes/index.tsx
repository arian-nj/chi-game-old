
import { createFileRoute, useRouter } from "@tanstack/react-router";
import { useState } from "react";
import { authBeforeLoad } from "../lib/auth";
import { MeComponent } from "../components/MeComponent";

export const Route = createFileRoute("/")({
	component: Index,
	beforeLoad: authBeforeLoad,
});
function Index() {
	const [selectedGame, setSelectedGame] = useState("");
	const router = useRouter();

	const games = ["3X3", "5X5"];

	function handlePlayClick() {
		router.navigate({ to: "/finder" });
	}

	return (
		<div className="relative h-screen w-screen bg-gradient-to-br from-blue-50 to-blue-100 overflow-hidden">
			<div className="flex justify-center gap-4 pt-16">
				{games.map((game) => (
					<GameChooser
						key={game}
						name={game}
						selected={selectedGame === game}
						onSelect={() => setSelectedGame(game)}
					/>
				))}
			</div>

			<div className="flex justify-center ">
				<MeComponent />
			</div>

			<button
				type="button"
				onClick={handlePlayClick}
				disabled={selectedGame === ""}
				className={`
					absolute bottom-0 left-0 w-full
					py-6 text-4xl font-bold rounded-none
					transition-all duration-300 ease-in-out
					focus:outline-none focus:ring-4 focus:ring-pink-300
					${selectedGame
						? "bg-gradient-to-r from-pink-500 to-orange-400 text-white hover:opacity-90 shadow-md"
						: "bg-gray-300 text-gray-400 cursor-not-allowed"}
				`}
			>
				Play
			</button>
		</div>
	);
}

type GameChooserProps = {
	name: string;
	selected: boolean;
	onSelect: () => void;
};

function GameChooser({ name, selected, onSelect }: GameChooserProps) {
	return (
		<button
			type="button"
			onClick={onSelect}
			className={`
				px-12 py-8 mx-3 rounded-xl text-2xl font-semibold
				backdrop-blur-md bg-orange-100/60 border
				transition-all duration-300 ease-in-out
				hover:scale-105 hover:shadow-md
				${selected ? 'border-orange-500 shadow-xl text-orange-800' : 'border-transparent text-gray-800'}
			`}
		>
			{name}
		</button>
	);
}


