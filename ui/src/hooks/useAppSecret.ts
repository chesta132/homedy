import { useState, useCallback } from "react";

/**
 * Manages app_secret prompt flow.
 * Secret is NEVER stored — user is asked every time a protected action is
 * triggered (backup, restore, config read/write).
 */
export function useAppSecret() {
  const [prompting, setPrompting] = useState(false);
  const [resolver, setResolver] = useState<((s: string | null) => void) | null>(null);

  /**
   * Opens the secret modal and returns a Promise that resolves with the
   * entered secret, or null if the user cancels.
   */
  const getSecret = useCallback((): Promise<string | null> => {
    return new Promise((resolve) => {
      setPrompting(true);
      // Store resolver via function form to avoid setState treating it as an updater
      setResolver(() => resolve);
    });
  }, []);

  const submitPrompt = useCallback(
    (s: string) => {
      setPrompting(false);
      resolver?.(s.trim() || null);
      setResolver(null);
    },
    [resolver]
  );

  const cancelPrompt = useCallback(() => {
    setPrompting(false);
    resolver?.(null);
    setResolver(null);
  }, [resolver]);

  return { prompting, getSecret, submitPrompt, cancelPrompt };
}
