import { useEffect, useState } from "react";
import { validateSession } from "../api/session";

export function useSessionValidation(sessionId?: string) {
  const [valid, setValid] = useState(false);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!sessionId) {
      setValid(false);
      setLoading(false);
      return;
    }

    validateSession(sessionId)
      .then(() => setValid(true))
      .catch(() => setValid(false))
      .finally(() => setLoading(false));
  }, [sessionId]);

  return { valid, loading };
}
