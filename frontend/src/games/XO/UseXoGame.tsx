import { useEffect, useState } from "react"
import { MoveActionType, type MoveActionData } from "./XoAction"
import type { SessionSocket } from "../../lib/SessionWs"


export function useXoGame(
	{ socketRef }
		: {
			socketRef: React.RefObject<SessionSocket>
		}) {

	const [squares, setSquares] = useState<(number)[]>(Array(9).fill(null));
	const [lastPlayed, setLastPlayed] = useState<number | null>(null);

	useEffect(() => {
		socketRef.current.HandleGameMessage = (gameAction) => {
			console.log(gameAction)
			switch (gameAction.action) {
				case MoveActionType:
					handleMoveAction(gameAction.adata)
			}
			console.log(gameAction)
		}
	}, [])

	const handleMoveAction = (adata: MoveActionData) => {
		const newSquares = squares.slice()
		newSquares[adata.index] = adata.value
		setSquares(newSquares)
		setLastPlayed(adata.index)
	}

	function handleClick(i: number) {
		console.log("clicked cell index ", i)
		// if (squares[i]) return;
		//
		// const nextSquares = squares.slice();
		// nextSquares[i] = xIsNext ? 1 : 2;
		// setSquares(nextSquares);
		// setLastPlayed(i);
	}



	return { squares, handleClick, lastPlayed }
}
