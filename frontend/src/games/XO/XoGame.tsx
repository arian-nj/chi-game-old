import { useXoGameOffline } from "./UseXoGameOffline";
import { motion, type Variants } from "motion/react"
//
// type XOGameProps = {
// 	socketRef: React.RefObject<SessionSocket>
// }

// export default function XoGame({ socketRef }: XOGameProps) {

const draw: Variants = {
	hidden: { pathLength: 0, opacity: 0 },
	visible: (i: number) => {
		const delay = i * 0.5
		return {
			pathLength: 1,
			opacity: 1,
			transition: {
				pathLength: { delay, type: "spring", duration: 1.5, bounce: 0 },
				opacity: { delay, duration: 0.01 },
			},
		}
	},
}
export default function XoGame() {
	const xSize = 100
	return (
		<motion.svg
			width={xSize}
			height={xSize}
			className={"bg-amber-400"}
			initial="hidden"
			animate="visible"
		>

			<motion.line
				x1={xSize - 10}
				y1={10}
				x2={10}
				y2={xSize - 10}
				stroke="#ed0920"
				variants={draw}
				custom={3}
				strokeWidth={20}
				strokeLinecap={"round"}
			/>
			<motion.line
				x1={10}
				y1={10}
				x2={xSize - 10}
				y2={xSize - 10}
				stroke="#ed0920"
				variants={draw}
				strokeWidth={20}
				custom={2.5}
				strokeLinecap={"round"}
			/>
		</motion.svg>
	)
	{/* <div className="flex flex-col items-center justify-center w-screen h-screen bg-green-200"> */ }
	{/* 	<Game2PlayersCard /> */ }
	{/* 	<Board socketRef={socketRef} /> */ }
	{/* 	<Board /> */ }
	{/* </div> */ }

	{/* <motion.circle */ }
	{/* 	cx={250} */ }
	{/* 	cy={250} */ }
	{/* 	r={100} */ }
	{/* 	stroke="#ff0088" */ }
	{/* 	strokeWidth={10} */ }
	{/* 	custom={1} */ }
	{/* 	variants={draw} */ }
	{/* 	strokeLinecap={"round"} */ }
	{/* /> */ }
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
