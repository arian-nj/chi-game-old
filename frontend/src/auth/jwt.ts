export async function ValidateInitdata(initData: string) {
  const res = await fetch("/api/auth/validate-init/", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ init_data: initData }),
  });
  if (!res.ok) throw new Error("Login failed");

  const data = await res.json();

  sessionStorage.setItem("jwt_token", data.token);
}

export async function ValidateDymmy(userId: number) {
  const res = await fetch("/api/auth/validate-dummy/", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ user_id: userId }),
  });
  if (!res.ok) throw new Error("Login failed");

  const data = await res.json();

  sessionStorage.setItem("jwt_token", data.token);
}

export function GetJwtToken(): string | null {
  return sessionStorage.getItem("jwt_token");
}
