import { useState } from "react";
import type { GameSessionSocket } from "../../lib/SessionWs";

const MoveActionType = "move"

type MoveActionData = {
	index: number
	value: number
}

type SquareProps = {
	value: string | null;
	onSquareClick: () => void;

	animate: boolean;
};

type XOProps = {
	socketRef: React.RefObject<GameSessionSocket>
}

export default function XoGame({ socketRef }: XOProps) {

	return (
		<div className="flex flex-col items-center justify-center w-screen h-screen bg-green-200">
			<Board socketRef={socketRef} />
		</div>
	);
}

type BoardProps = {
	socketRef: React.RefObject<GameSessionSocket>
}

function Board({ socketRef }: BoardProps) {
	const handleMoveAction = (adata: MoveActionData) => {
		const newSquares = squares.slice()
		newSquares[adata.index] = adata.value
		setSquares(newSquares)
		setLastPlayed(adata.index)
	}

	socketRef.current.HandleGameAction = (gameAction) => {
		console.log(gameAction)
		switch (gameAction.action) {
			case MoveActionType:
				handleMoveAction(gameAction.adata)
		}
		console.log(gameAction)
	}
	const [xIsNext, setXIsNext] = useState(true);
	const [squares, setSquares] = useState<(number)[]>(Array(9).fill(null));
	const [lastPlayed, setLastPlayed] = useState<number | null>(null);

	function handleClick(i: number) {
		if (squares[i]) return;

		const nextSquares = squares.slice();
		nextSquares[i] = xIsNext ? 1 : 2;
		setSquares(nextSquares);
		setXIsNext(!xIsNext);
		setLastPlayed(i);
	}

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
