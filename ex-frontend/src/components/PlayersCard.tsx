import type { Account } from "../gen/account/v1/account_pb"
import { useEffect } from "react"
import toast from "react-hot-toast"
import { useQuery as connectQuery } from "@connectrpc/connect-query"
import { getSessionOpponent } from "../gen/session/v1/session-SessionService_connectquery"
import { getMe } from "../gen/account/v1/account-AccountService_connectquery"

import { motion } from "motion/react"

export function Game2PlayersCard() {
	const { isPending: isOppPending, error: oppError, data: oppData } = connectQuery(getSessionOpponent);
	const { data: meData } = connectQuery(getMe)

	useEffect(() => {
		if (oppError) {
			toast.error(`Can't get opponent player: ${oppError.message}`);
			console.error("can't get opp player", oppError);
		}
	}, [oppError]);

	if (isOppPending) {
		return <div>Loading opponent...</div>;
	}

	return (
		<>
			<motion.div
				initial={{ y: -100 }}
				animate={{ y: 0 }}
				transition={{ duration: 1 }}


				className="absolute top-0 left-1/2 -translate-x-1/2 
             flex justify-between gap-5 px-10 py-4
             bg-[#4B5043] border-b-4 border-x-4 border-[#2E3228]
             rounded-b-2xl shadow-xl
             max-w-lg w-[90%] z-50"

			>
				{meData && meData.account &&
					<PlayerCard account={meData.account} isActive={true} />
				}

				{oppData && oppData.opponent && (
					<PlayerCard account={oppData.opponent} isActive={false} />
				)}

			</motion.div>
		</>
	);
}

type PlayerCardProps = {
	account: Account
	isActive: boolean
}

function PlayerCard({ account, isActive }: PlayerCardProps) {
	return (
		<div className={`text-white p-4 text-center rounded-2xl min-w-24
						${isActive ? "bg-[#4283AB]/100" : "bg-[#4283AB]/30"}
						`}>
			<h3 className="text-xl font-semibold mb-1">{account.name}</h3>
			<p className="text-sm opacity-80">ID: {account.id}</p>
		</div >
	);
}

