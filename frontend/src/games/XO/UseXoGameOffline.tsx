import { useState } from "react"
export function useXoGameOffline() {

	const [squares, setSquares] = useState<(number)[]>(Array(9).fill(null));
	const [lastPlayed, setLastPlayed] = useState<number | null>(null);

	const DoMove = (index: number, value: number) => {
		setSquares(prevSquares => {
			const newSquares = [...prevSquares]
			newSquares[index] = value
			return newSquares
		})
		setLastPlayed(index)
	}

	function handleClick(i: number) {
		console.log("clicked cell index ", i)
		DoMove(i, 1)
		// if (squares[i]) return;
		//
		// const nextSquares = squares.slice();
		// nextSquares[i] = xIsNext ? 1 : 2;
		// setSquares(nextSquares);
		// setLastPlayed(i);
	}



	return { squares, handleClick, lastPlayed }
}
