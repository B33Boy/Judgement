import { useEffect, useState } from "react";
import type { NavigateFunction } from "react-router-dom";

const API_BASE = `http://localhost:${import.meta.env.VITE_PORT}/api`;

interface SessionValidationResult {
  valid: boolean;
  loading: boolean;
}

export function useSessionValidation(
  sessionId: string | undefined,
  navigate: NavigateFunction
): SessionValidationResult {
  const [valid, setValid] = useState(false);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!sessionId) {
      navigate("/");
      return;
    }

    async function validate() {
      try {
        const res = await fetch(API_BASE + `/session/${sessionId}`);

        if (!res.ok) {
          navigate("/");
          return;
        }

        setValid(true);
      } finally {
        setLoading(false);
      }
    }

    validate();
  }, [sessionId, navigate]);

  return { valid, loading };
}
