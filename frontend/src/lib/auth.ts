import { redirect } from "@tanstack/react-router";
import { GetBaseUrl } from "./baseURL";

export function authBeforeLoad() {
	const jwtToken = GetJwtToken()
	if (jwtToken == null) {
		throw redirect({
			to: "/login",
			search: {
				redirect: location.pathname + location.search
			}
		})
	}
}

export function GetJwtToken(): string | null {
	return sessionStorage.getItem("jwt_token");
}

export function SetJwtToken(token: string) {
	sessionStorage.setItem("jwt_token", token);
}

class MeClass {
	ID: number
	constructor(id: number) {
		this.ID = id
	}
}
function fetchWithAuth(input: RequestInfo, init?: RequestInit) {
	const token = GetJwtToken();

	const headers = new Headers(init?.headers);
	if (token) {
		headers.set("Authorization", `Bearer ${token}`);
	}

	return fetch(input, {
		...init,
		headers,
	});
}
export async function GetMe(): Promise<MeClass> {
	const response = await fetchWithAuth(GetBaseUrl() + "/api/auth/me")
	if (response.status != 200) {
		throw new Error(response.statusText)
	}
	const json_data = await response.json()
	console.log("data ", json_data)
	return new MeClass(json_data.id)
}
