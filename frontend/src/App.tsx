import "./App.css";
import { Routes, Route, Navigate } from "react-router-dom";
import { GameProvider } from "./context/GameContext";

import HomePage from "./pages/Home";
import SessionPage from "./pages/Session";
import NotFoundPage from "./pages/NotFound";
import GamePage from "./pages/Game";

function App() {
  return (
    <GameProvider>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/session" element={<Navigate to="/" replace />} />
        <Route path="/session/:sessionId" element={<SessionPage />} />
        <Route path="/game/:sessionId/:playerName" element={<GamePage />} />

        <Route path="*" element={<NotFoundPage />} />
      </Routes>
    </GameProvider>
  );
}
export default App;
