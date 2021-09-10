export function getAccessToken(): string | null {
  return localStorage.getItem("access_token");
}

export function setAccessToken(token: string): void {
  return localStorage.setItem("access_token", token);
}
