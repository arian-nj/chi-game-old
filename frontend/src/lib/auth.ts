import { redirect } from "@tanstack/react-router";

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


