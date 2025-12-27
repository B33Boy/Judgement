import "./App.css";

import { Routes, Route, Navigate } from "react-router-dom";

import HomePage from "./pages/Home";
import SessionPage from "./pages/Session";
import NotFoundPage from "./pages/NotFound";

function App() {
  return (
    <Routes>
      <Route path="/" element={<HomePage />} />
      <Route path="/session" element={<Navigate to="/" replace />} />
      <Route path="/session/:sessionId" element={<SessionPage />} />

      <Route path="*" element={<NotFoundPage />} />
    </Routes>
  );
}

export default App;
