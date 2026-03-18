import { ServerError } from "@/services/server/ServerError";
import type { StateErrorServer } from "@/types/server";

/**
 * Maps any caught error to the global error state.
 */
export const handleError = (
  err: unknown,
  setError: React.Dispatch<React.SetStateAction<StateErrorServer | null>>
) => {
  if (err instanceof ServerError) {
    setError({ code: err.getCode(), message: err.getMessage(), title: err.getTitle() });
  } else if (err instanceof Error) {
    if (err.message.toLowerCase().includes("network")) {
      setError({ code: "BAD_GATEWAY", message: "Unable to connect to server. Check your connection.", title: "Connection Error" });
    } else {
      setError({ code: "SERVER_ERROR", message: err.message, title: "Something went wrong" });
    }
  } else {
    setError({ code: "SERVER_ERROR", message: "An unexpected error occurred.", title: "Something went wrong" });
  }
};

/**
 * Like handleError but also sets field-level errors for form inputs.
 */
export const handleFormError = <T extends Record<string, string>>(
  err: unknown,
  setFormError: React.Dispatch<React.SetStateAction<T>>,
  setError: React.Dispatch<React.SetStateAction<StateErrorServer | null>>
) => {
  if (err instanceof ServerError) {
    const field = err.getField();
    if (field) {
      setFormError((prev) => ({ ...prev, [field]: err.getMessage() }));
      return;
    }
  }
  handleError(err, setError);
};
