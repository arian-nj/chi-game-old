import { useQuery, type QueryFunctionContext } from "@tanstack/react-query";
import GameSession from "./GameSession";
import { useEffect, useState } from "react";
import { GetJwtToken } from "./lib/auth";

export default function App() {
	const [token, setToken] = useState<string | null>(GetJwtToken());
	if (token === null) return <Auth onAuthSuccess={setToken} />;

	return token == null ? <Auth onAuthSuccess={setToken} /> : <GameSession />;
};
type AuthProps = {
	onAuthSuccess: (token: string) => void;
};
function Auth({ onAuthSuccess }: AuthProps) {
	const [userId, setUserId] = useState<number | null>(null);

	useEffect(() => {
		const userIDStr = prompt("What is your user ID", "1");
		const userId = Number(userIDStr);
		if (!Number.isNaN(userId)) setUserId(userId);
	}, []);

	const { isPending, error, data } = useQuery({
		enabled: userId !== null, // only run if userId is set
		queryKey: ["dummy_auth", userId!],
		queryFn: ValidateDummy,
	});

	useEffect(() => {
		if (data) {
			sessionStorage.setItem("jwt_token", data.token); // use sessionStorage consistently
			onAuthSuccess(data.token)
		}
	}, [data]);

	if (isPending) return <p>Loading...</p>;
	if (error) return <p>Error: {(error as Error).message}</p>;

	// `data` already handled in useEffect
	return <p>Validated</p>;
}
type AuthResponse = { token: string };

async function ValidateDummy({ queryKey }: QueryFunctionContext<[string, number]>): Promise<AuthResponse> {
	const [, userId] = queryKey
	const response = await fetch("/api/auth/validate-dummy/", {
		method: "POST",
		headers: { "Content-Type": "application/json" },
		body: JSON.stringify({ user_id: userId }),
	});
	if (!response.ok) throw new Error("Login failed");

	return response.json();
}

// export async function ValidateInitdata(initData: string) {
// 	const res = await fetch("/api/auth/validate-init/", {
// 		method: "POST",
// 		headers: { "Content-Type": "application/json" },
// 		body: JSON.stringify({ init_data: initData }),
// 	});
// 	if (!res.ok) throw new Error("Login failed");
//
// 	const data = await res.json();
//
// 	sessionStorage.setItem("jwt_token", data.token);
// }
