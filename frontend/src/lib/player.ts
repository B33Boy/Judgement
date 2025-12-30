export function getPlayerName(): string {
  return localStorage.getItem("playerName") ?? "";
}

export function setPlayerName(name: string) {
  localStorage.setItem("playerName", name);
}
