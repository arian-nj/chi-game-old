import type { SessionSocket } from "../../lib/SessionWs";
import { Board } from "./XoBoard";

type XOGameProps = {
	socketRef: React.RefObject<SessionSocket>
}

export default function XoGame({ socketRef }: XOGameProps) {

	return (
		<div className="flex flex-col items-center justify-center w-screen h-screen bg-green-200">
			<Board socketRef={socketRef} />
		</div>
	);
}


