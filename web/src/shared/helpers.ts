export function getAccessToken(): string | null {
  return localStorage.getItem("access_token");
}

export function setAccessToken(token: string): void {
  return localStorage.setItem("access_token", token);
}

export function getAuthorization(): string | null {
  const token = localStorage.getItem("access_token");
  if (!token) {
    return null;
  }
  return `Bearer ${token}`;
}
