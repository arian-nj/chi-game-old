import type { SessionSocket } from "../../lib/SessionWs"
import { useXoGame } from "./UseXoGame"

type BoardProps = {
	socketRef: React.RefObject<SessionSocket>
}
export function Board({ socketRef }: BoardProps) {
	const { squares, handleClick, lastPlayed } = useXoGame({ socketRef })
	return (
		<div
			className="grid grid-cols-3 gap-1">
			{squares.map((value, i) => (
				<Square
					key={i}
					value={value == 1 ? "X" : value == 2 ? "O" : ""}
					animate={lastPlayed === i}
					onSquareClick={() => handleClick(i)} />
			))}
		</div>
	);
}
type SquareProps = {
	value: string | null;
	onSquareClick: () => void;
	animate: boolean;
};


function Square({ value, onSquareClick, animate }: SquareProps) {
	return (
		<button
			onClick={onSquareClick}
			className="w-20 h-20 bg-amber-50 hover:bg-gray-200
			flex items-center justify-center border-2 border-gray-700
			text-red-700 text-5xl font-bold
			transition-all duration-200 ease-out
			"
		>
			{value && (
				<span className={animate ? "animate-pop" : ""}>
					{value}
				</span>
			)}
		</button>
	);
}
