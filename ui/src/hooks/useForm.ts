import { useState } from "react";
import { FormValidator } from "@/services/form-validator/FormValidator";
import type { FormConfig, FormFields } from "@/types/form";

export type FormGroup<T extends FormFields> = {
  readonly form: [T, React.Dispatch<React.SetStateAction<T>>];
  readonly error: [Record<keyof T, string>, React.Dispatch<React.SetStateAction<Record<keyof T, string>>>];
  readonly validate: {
    validateForm: () => boolean;
    validateField: (field: Partial<T>) => boolean;
  };
  readonly resetForm: () => void;
};

/**
 * Form state and validation hook.
 * Ported from hoshify-client useForm pattern.
 */
const useForm = <T extends FormFields>(schema: T, config: FormConfig<T>): FormGroup<T> => {
  const [form, setForm] = useState<T>(schema);
  const [formError, setFormError] = useState<Record<keyof T, string>>(
    Object.fromEntries(Object.keys(schema).map((k) => [k, ""])) as Record<keyof T, string>
  );

  const validator = new FormValidator(form, config);

  const validateForm = () => {
    const { errors, hasError } = validator.validateForm();
    if (hasError) {
      setFormError((prev) => ({ ...prev, ...errors }));
    } else {
      setFormError(Object.fromEntries(Object.keys(formError).map((k) => [k, ""])) as Record<keyof T, string>);
    }
    return !hasError;
  };

  const validateField = (field: Partial<T>) => {
    const fieldValidator = new FormValidator({ ...form, ...field }, config);
    const { errors, hasError } = fieldValidator.validateFields(field);
    setForm((prev) => ({ ...prev, ...field }));
    if (hasError) {
      setFormError((prev) => ({ ...prev, ...errors }));
    } else {
      const cleared = Object.fromEntries(Object.keys(field).map((k) => [k, ""])) as Partial<Record<keyof T, string>>;
      setFormError((prev) => ({ ...prev, ...cleared }));
    }
    return !hasError;
  };

  const resetForm = () => {
    setForm(schema);
    setFormError(Object.fromEntries(Object.keys(schema).map((k) => [k, ""])) as Record<keyof T, string>);
  };

  return {
    form: [form, setForm],
    error: [formError, setFormError],
    validate: { validateForm, validateField },
    resetForm,
  };
};

export default useForm;
