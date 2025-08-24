import { useEffect, useState } from "react"
import type { SessionSocket } from "../../lib/SessionWs"
import * as XoBuff from "../../gen/xo_game/v1/xo_pb";
import { create, toBinary } from "@bufbuild/protobuf";
import { SessionMessageSchema } from "../../gen/session/v1/session_pb";


export function useXoGame(
	{ socketRef }
		: {
			socketRef: React.RefObject<SessionSocket>
		}) {

	const [squares, setSquares] = useState<(number)[]>(Array(9).fill(null));
	const [lastPlayed, setLastPlayed] = useState<number | null>(null);

	useEffect(() => {
		socketRef.current.HandleGameMessage = (gameMessage) => {
			console.log(gameMessage)
			if (gameMessage.game.case != "xo") {
				throw Error("non xo message ended up in xo")
			}
			const payload = gameMessage.game.value.payload
			switch (payload.case) {
				case "move":
					handleMoveAction(payload.value)
					break

				case "playResponse":
					handlePlayResponse(payload.value)
					break
			}
		}
	}, [])

	const DoMove = (index: number, value: number) => {
		setSquares(prevSquares => {
			const newSquares = [...prevSquares]
			newSquares[index] = value
			return newSquares
		})
		setLastPlayed(index)
	}
	const handleMoveAction = (moveData: XoBuff.Move) => {
		DoMove(moveData.cellIndex, 1)
	}

	const handlePlayResponse = (playResponse: XoBuff.PlayResponse) => {
		if (playResponse.isValid) {
			if (playResponse.play) {
				DoMove(playResponse.play.cellIndex, 1)
			}
		} else {
			alert("wrong move")
		}
	}



	function handleClick(i: number) {
		console.log("clicked cell index ", i)
		const newSessionMsg = create(SessionMessageSchema, {
			content: {
				case: "game",
				value: {
					game: {
						case: "xo",
						value: { payload: { case: "play", value: { cellIndex: i } } }
					}
				}
			}
		})
		const bytes = toBinary(SessionMessageSchema, newSessionMsg)
		socketRef.current.send(bytes)
		// if (squares[i]) return;
		//
		// const nextSquares = squares.slice();
		// nextSquares[i] = xIsNext ? 1 : 2;
		// setSquares(nextSquares);
		// setLastPlayed(i);
	}



	return { squares, handleClick, lastPlayed }
}
