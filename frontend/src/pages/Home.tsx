import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { createSession, validateSession } from "../api/session";

function CreateRoomForm({ playerName }: { playerName: string }) {
  const navigate = useNavigate();

  async function handleSubmit() {
    if (!playerName.trim()) {
      alert("Enter a player name");
      return;
    }
    localStorage.setItem("playerName", playerName);

    const sessionId = await createSession();

    navigate(`/session/${sessionId}`);
  }

  return (
    <button onClick={handleSubmit} type="button">
      Create Room
    </button>
  );
}

function JoinRoomForm({ playerName }: { playerName: string }) {
  const navigate = useNavigate();
  const [joinCode, setJoinCode] = useState("");

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();

    if (!playerName.trim()) {
      alert("Enter a player name");
      return;
    }
    localStorage.setItem("playerName", playerName);

    if (!joinCode.trim()) return;

    try {
      await validateSession(joinCode);
      navigate(`/session/${joinCode}`);
    } catch {
      alert("Invalid Session Id");
    }
  }

  return (
    <form onSubmit={handleSubmit}>
      <input
        value={joinCode}
        placeholder="Room Code"
        onChange={(e) => setJoinCode(e.target.value)}
        type="text"
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

      <div className="card">
        <input
          type="text"
          placeholder="Player Name"
          value={playerName}
          onChange={(e) => setPlayerName(e.target.value)}
        />
        <CreateRoomForm playerName={playerName} />
        <JoinRoomForm playerName={playerName} />
      </div>
    </>
  );
}
