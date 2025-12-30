const API_BASE = `http://localhost:${import.meta.env.VITE_PORT}/api`;

export async function createSession(): Promise<string> {
  const res = await fetch(`${API_BASE}/session`, { method: "POST" });

  if (!res.ok) {
    throw new Error("Failed to create session");
  }

  const data = await res.json();
  return data.sessionId;
}

export async function validateSession(sessionId: string): Promise<void> {
  const res = await fetch(`${API_BASE}/session/${sessionId}`);

  if (!res.ok) {
    throw new Error("Invalid session");
  }
}
