import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { createSession, validateSession } from "../api/session";
import { setPlayerName as savePlayerName } from "../lib/player";

export default function HomePage() {
  // States
  const [playerName, setPlayerName] = useState("");
  const [joinCode, setJoinCode] = useState("");
  const navigate = useNavigate();

  const validateInput = () => {
    if (!playerName.trim()) {
      alert("Please enter a name");
      return false;
    }
    savePlayerName(playerName);
    return true;
  };

  const handleCreate = async () => {
    if (!validateInput()) return;
    try {
      const sessionId = await createSession();
      navigate(`/session/${sessionId}`);
    } catch (e) {
      alert("Failed to create session");
    }
  };

  const handleJoin = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!validateInput() || !joinCode.trim()) return;

    try {
      await validateSession(joinCode);
      navigate(`/session/${joinCode}`);
    } catch {
      alert("Invalid session code");
    }
  };

  return (
    <div className="home-container">
      <h1>Judgement</h1>
      <input
        placeholder="Player Name"
        value={playerName}
        onChange={(e) => setPlayerName(e.target.value)}
        maxLength={20}
      />

      <div className="actions">
        <button onClick={handleCreate}>Create Room</button>

        <form onSubmit={handleJoin} style={{ marginTop: "20px" }}>
          <input
            placeholder="Room Code"
            value={joinCode}
            onChange={(e) => setJoinCode(e.target.value)}
            maxLength={8}
          />
          <button type="submit">Join</button>
        </form>
      </div>
    </div>
  );
}
