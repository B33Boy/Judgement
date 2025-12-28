import { useNavigate } from "react-router-dom";

const API_BASE = `http://localhost:${import.meta.env.VITE_PORT}/api`;

function CreateRoomButton() {
  const navigate = useNavigate();

  async function handleSubmit() {
    try {
      const res = await fetch(API_BASE + "/session", {
        method: "POST",
      });

      if (!res.ok) {
        throw new Error("Failed to create session");
      }

      const data = await res.json();
      const sessionId = data.sessionId;

      navigate(`/session/${sessionId}`);
    } catch (err) {
      alert("Error creating room");
      console.error(err);
    }
  }

  return (
    <button onClick={handleSubmit} type="button">
      Create Room
    </button>
  );
}

function JoinRoomForm() {
  const navigate = useNavigate();

  async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault(); // STOP native submission

    const formData = new FormData(e.currentTarget);
    const joinCode = formData.get("joinCode")?.toString();

    if (!joinCode) return;

    try {
      const res = await fetch(
        `http://localhost:${import.meta.env.VITE_PORT}/api/session/${joinCode}`
      );

      if (!res.ok) {
        alert("Invalid Session Id");
        return;
      }

      navigate(`/session/${joinCode}`);
    } catch (err) {
      alert("Error joining room");
      console.error(err);
    }
  }

  return (
    <form onSubmit={handleSubmit}>
      <input name="joinCode" type="text" />
      <button type="submit">Join</button>
    </form>
  );
}

export default function HomePage() {
  return (
    <>
      <h1>Judgement</h1>
      <div className="card">
        <CreateRoomButton></CreateRoomButton>
        <JoinRoomForm></JoinRoomForm>
      </div>
    </>
  );
}
