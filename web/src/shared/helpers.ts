export function getAccessToken(): string | undefined {
  return localStorage.get("access_token");
}

export function setAccessToken(token: string): void {
  return localStorage.set("access_token", token);
}
