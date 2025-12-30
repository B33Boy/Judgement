import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { createSession, validateSession } from "../api/session";
import { setPlayerName } from "../lib/player";

function requirePlayerName(playerName: string): boolean {
  if (!playerName.trim()) {
    alert("Enter a player name");
    return false;
  }
  setPlayerName(playerName);
  return true;
}

function CreateRoomForm({ playerName }: { playerName: string }) {
  const navigate = useNavigate();

  async function handleClick() {
    if (!requirePlayerName(playerName)) return;

    const sessionId = await createSession();
    navigate(`/session/${sessionId}`);
  }

  return <button onClick={handleClick}>Create Room</button>;
}

function JoinRoomForm({ playerName }: { playerName: string }) {
  const navigate = useNavigate();
  const [joinCode, setJoinCode] = useState("");

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();

    if (!requirePlayerName(playerName)) return;
    if (!joinCode.trim()) return;

    try {
      await validateSession(joinCode);
      navigate(`/session/${joinCode}`);
    } catch {
      alert("Invalid session");
    }
  }

  return (
    <form onSubmit={handleSubmit}>
      <input
        placeholder="Room Code"
        value={joinCode}
        onChange={(e) => setJoinCode(e.target.value)}
      />
      <button type="submit">Join</button>
    </form>
  );
}

export default function HomePage() {
  const [playerName, setPlayerName] = useState("");

  return (
    <>
      <h1>Judgement</h1>
      <input
        placeholder="Player Name"
        value={playerName}
        onChange={(e) => setPlayerName(e.target.value)}
      />
      <CreateRoomForm playerName={playerName} />
      <JoinRoomForm playerName={playerName} />
    </>
  );
}
