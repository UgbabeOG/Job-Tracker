
import { useState } from "react";
import { Button, FormField } from "../components/ui";
import type { Application, ApplicationStatus } from "./types";

type ApplicationFormData = Omit<
  Application,
  "id" | "created_at" | "updated_at"
>;

interface ApplicationFormProps {
  initialApplication?: Application | null;
  onSubmit: (data: ApplicationFormData) => Promise<void>;
  onCancel?: () => void;
  submitting?: boolean;
}

export default function ApplicationForm({
  initialApplication,
  onSubmit,
  onCancel,
  submitting = false,
}: ApplicationFormProps) {
  const [company, setCompany] = useState(
    initialApplication?.company ?? ""
  );
  const [jobTitle, setJobTitle] = useState(
    initialApplication?.job_title ?? ""
  );
  const [jobUrl, setJobUrl] = useState(
    initialApplication?.job_url ?? ""
  );
  const [dateApplied, setDateApplied] = useState(
    initialApplication?.date_applied ?? ""
  );
  const [deadline, setDeadline] = useState(
    initialApplication?.deadline ?? ""
  );
  const [status, setStatus] = useState<ApplicationStatus>(
    initialApplication?.status ?? "Saved"
  );
  const [notes, setNotes] = useState(
    initialApplication?.notes ?? ""
  );
  const [error, setError] = useState("");

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError("");

    if (!company.trim() || !jobTitle.trim()) {
      setError("Company name and job title are required.");
      return;
    }

    if (jobUrl.trim()) {
      try {
        new URL(jobUrl);
      } catch {
        setError("Please enter a valid job URL, including https://.");
        return;
      }
    }

    const applicationData: ApplicationFormData = {
      company: company.trim(),
      job_title: jobTitle.trim(),
      job_url: jobUrl.trim() || null,
      date_applied: dateApplied || null,
      deadline: deadline || null,
      status,
      notes: notes.trim() || null,
    };

    try {
      await onSubmit(applicationData);
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : "Something went wrong. Please try again."
      );
    }
  }

  return (
    <form className="application-form" onSubmit={handleSubmit}>
      <div className="application-form-grid">
        <FormField
          id="company"
          label="Company name *"
          value={company}
          onChange={(event) => setCompany(event.target.value)}
          placeholder="e.g. Google"
          required
        />

        <FormField
          id="job-title"
          label="Job title *"
          value={jobTitle}
          onChange={(event) => setJobTitle(event.target.value)}
          placeholder="e.g. Frontend Developer"
          required
        />

        <FormField
          id="job-url"
          label="Job URL"
          type="url"
          value={jobUrl}
          onChange={(event) => setJobUrl(event.target.value)}
          placeholder="https://example.com/job"
        />

        <FormField
          id="date-applied"
          label="Date applied"
          type="date"
          value={dateApplied}
          onChange={(event) => setDateApplied(event.target.value)}
        />

        <FormField
          id="deadline"
          label="Application deadline"
          type="date"
          value={deadline}
          onChange={(event) => setDeadline(event.target.value)}
        />

        <div className="application-form-field">
          <label htmlFor="status">Status</label>
          <select
            id="status"
            value={status}
            onChange={(event) =>
              setStatus(event.target.value as ApplicationStatus)
            }
          >
            <option value="Saved">Saved</option>
            <option value="Applied">Applied</option>
            <option value="Interview">Interview</option>
            <option value="Offer">Offer</option>
            <option value="Rejected">Rejected</option>
          </select>
        </div>
      </div>

      <div className="application-form-field">
        <label htmlFor="notes">Notes</label>
        <textarea
          id="notes"
          value={notes}
          onChange={(event) => setNotes(event.target.value)}
          placeholder="Add interview details, requirements, or reminders..."
          rows={4}
        />
      </div>

      {error && <p className="application-form-error">{error}</p>}

      <div className="application-form-actions">
        {onCancel && (
          <Button
            type="button"
            variant="secondary"
            onClick={onCancel}
            disabled={submitting}
          >
            Cancel
          </Button>
        )}

        <Button type="submit" disabled={submitting}>
          {submitting
            ? "Saving..."
            : initialApplication
              ? "Update Application"
              : "Save Application"}
        </Button>
      </div>
    </form>
  );
}
