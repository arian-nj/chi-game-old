import { useEffect, useRef, useState } from "react";
import { useXoGameOffline } from "./UseXoGameOffline";
import { Application, extend } from "@pixi/react"
import { Assets, Container, Graphics, Sprite, Texture } from "pixi.js";
//
// type XOGameProps = {
// 	socketRef: React.RefObject<SessionSocket>
// }

// export default function XoGame({ socketRef }: XOGameProps) {
// extend tells @pixi/react what Pixi.js components are available
extend({
	Container,
	Graphics,
	Sprite,
});

export default function XoGame() {
	return (

		<div className="flex flex-col items-center justify-center w-screen h-screen bg-green-200">
			{/* <Game2PlayersCard /> */}
			{/* <Board socketRef={socketRef} /> */}
			<Application>
				<BunnySprite />
			</Application>

			<Board />
		</div>
	)
}

function BunnySprite() {
	const spriteRef = useRef(null)

	const [texture, setTexture] = useState(Texture.EMPTY)
	const [isActive, setIsActive] = useState(false)

	useEffect(() => {
		if (texture === Texture.EMPTY) {
			Assets.load('https://pixijs.com/assets/bunny.png')
				.then((result) => {
					setTexture(result)
				});
		}
	}, [texture])

	return (
		<pixiSprite
			ref={spriteRef}
			anchor={0.5}
			eventMode={'static'}
			scale={isActive ? 1 : 1.5}
			onClick={() => setIsActive(!isActive)}
			texture={texture}
			x={200}
			y={200}
		/>
	)
}

// type BoardProps = {
// 	socketRef: React.RefObject<SessionSocket>
// }
// export function Board({ socketRef }: BoardProps) {
// const { squares, handleClick, lastPlayed } = useXoGameOnline({ socketRef })
export function Board() {
	const { squares, handleClick, lastPlayed } = useXoGameOffline()

	return (
		<>

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
		</>
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
			className="w-20 h-20 hover:bg-gray-200
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
