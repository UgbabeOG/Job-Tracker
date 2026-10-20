import type {
  ButtonHTMLAttributes,
  InputHTMLAttributes,
  ReactNode,
} from "react";

type ButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: "primary" | "secondary" | "quiet";
};

export function Button({
  className = "",
  variant = "primary",
  ...props
}: ButtonProps) {
  return (
    <button
      className={`button button-${variant} ${className}`.trim()}
      {...props}
    />
  );
}

type FormFieldProps = Omit<InputHTMLAttributes<HTMLInputElement>, "id"> & {
  label: string;
  id: string;
  hint?: string;
  error?: string;
};

export function FormField({
  label,
  id,
  hint,
  error,
  "aria-describedby": existingDescription,
  ...inputProps
}: FormFieldProps) {
  const descriptionId = `${id}-description`;
  const errorId = `${id}-error`;
  const describedBy = [
    existingDescription,
    hint ? descriptionId : undefined,
    error ? errorId : undefined,
  ]
    .filter(Boolean)
    .join(" ");

  return (
    <div className="form-field">
      <label htmlFor={id}>{label}</label>
      <input
        id={id}
        aria-describedby={describedBy || undefined}
        aria-invalid={error ? true : undefined}
        {...inputProps}
      />
      {hint && (
        <small className="field-hint" id={descriptionId}>
          {hint}
        </small>
      )}
      {error && (
        <small className="field-error" id={errorId} role="alert">
          {error}
        </small>
      )}
    </div>
  );
}

export type ApplicationStatus =
  | "Saved"
  | "Applied"
  | "Interview"
  | "Offer"
  | "Rejected";

const statusClass: Record<ApplicationStatus, string> = {
  Saved: "status-saved",
  Applied: "status-applied",
  Interview: "status-interview",
  Offer: "status-offer",
  Rejected: "status-rejected",
};

export function StatusBadge({ status }: { status: ApplicationStatus }) {
  return (
    <span className={`status-badge ${statusClass[status]}`}>{status}</span>
  );
}

export function LoadingState({ label = "Loading..." }: { label?: string }) {
  return (
    <p className="state-message" role="status" aria-live="polite">
      <span className="loading-indicator" aria-hidden="true" />
      {label}
    </p>
  );
}

export function ErrorMessage({ message }: { message: string }) {
  return (
    <p className="error-message" role="alert">
      {message}
    </p>
  );
}

type EmptyStateProps = {
  title: string;
  description: string;
  action?: ReactNode;
};

export function EmptyState({
  title,
  description,
  action,
}: EmptyStateProps) {
  return (
    <section className="empty-state">
      <h2>{title}</h2>
      <p>{description}</p>
      {action}
    </section>
  );
}
