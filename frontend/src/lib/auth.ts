
export function GetJwtToken(): string | null {
  return sessionStorage.getItem("jwt_token");
}

export function SetJwtToken(token: string) {
  sessionStorage.setItem("jwt_token", token);
}
