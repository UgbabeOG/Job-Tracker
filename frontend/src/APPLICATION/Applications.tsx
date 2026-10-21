
import { useCallback, useEffect, useState } from "react";
import {
  Button,
  StatusBadge,
  LoadingState,
  ErrorMessage,
  EmptyState,
} from "../components/ui";

// Applications Management — Issue #7
import type { Application, ApplicationStatus } from "./types";
import {
  getApplications,
  createApplication,
  updateApplication,
  deleteApplication,
} from "./Application";
import ApplicationForm from "./ApplicationForm";
import "./Application.css";

export default function Applications() {
  const [applications, setApplications] = useState<Application[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [showForm, setShowForm] = useState(false);
  const [editingApplication, setEditingApplication] =
    useState<Application | null>(null);
  const [saving, setSaving] = useState(false);
  const [deletingId, setDeletingId] = useState<number | null>(null);

  // Load application records from the API.
  const loadApplications = useCallback(async () => {
    setLoading(true);
    setError("");

    try {
      const data = await getApplications();
      setApplications(data);
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : "Failed to load applications."
      );
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void loadApplications();
  }, [loadApplications]);

  // Create a new application or save changes to an existing one.
  async function handleSave(
    data: Omit<Application, "id" | "created_at" | "updated_at">
  ) {
    setSaving(true);
    setError("");

    try {
      if (editingApplication) {
        const updated = await updateApplication(
          editingApplication.id,
          data
        );

        setApplications((current) =>
          current.map((app) =>
            app.id === updated.id ? updated : app
          )
        );
      } else {
        const created = await createApplication(data);
        setApplications((current) => [created, ...current]);
      }

      setShowForm(false);
      setEditingApplication(null);
    } finally {
      setSaving(false);
    }
  }

  // Update an application's status.
  async function handleStatusChange(
    application: Application,
    status: ApplicationStatus
  ) {
    setError("");

    try {
      const updated = await updateApplication(application.id, {
        status,
      });

      setApplications((current) =>
        current.map((app) =>
          app.id === updated.id ? updated : app
        )
      );
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : "Failed to update application status."
      );
    }
  }

  // Confirm before deleting an application.
  async function handleDelete(application: Application) {
    const confirmed = window.confirm(
      `Are you sure you want to delete the application for ${application.company}?`
    );

    if (!confirmed) return;

    setDeletingId(application.id);
    setError("");

    try {
      await deleteApplication(application.id);

      setApplications((current) =>
        current.filter((app) => app.id !== application.id)
      );
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : "Failed to delete application."
      );
    } finally {
      setDeletingId(null);
    }
  }

  function handleEdit(application: Application) {
    setEditingApplication(application);
    setShowForm(true);
  }

  function handleAdd() {
    setEditingApplication(null);
    setShowForm(true);
  }

  function handleCancel() {
    setShowForm(false);
    setEditingApplication(null);
  }

  return (
    <section className="applications-page">
      <header className="applications-header">
        <div>
          <h1>My Applications</h1>
          <p>Track and manage your job search in one place.</p>
        </div>

        {!showForm && (
          <Button onClick={handleAdd}>
            + Add Application
          </Button>
        )}
      </header>

      {error && <ErrorMessage message={error} />}

      {showForm && (
        <div className="application-form-panel">
          <h2>
            {editingApplication
              ? "Edit Application"
              : "Add New Application"}
          </h2>

          <ApplicationForm
            key={editingApplication?.id ?? "new"}
            initialApplication={editingApplication}
            onSubmit={handleSave}
            onCancel={handleCancel}
            submitting={saving}
          />
        </div>
      )}

      <div className="applications-summary">
        <span>
          {applications.length}{" "}
          {applications.length === 1
            ? "application"
            : "applications"}
        </span>
      </div>

      {loading ? (
        <LoadingState />
      ) : applications.length === 0 ? (
        <EmptyState
          title="No applications yet"
          description="Add your first job application to start tracking your progress."
        />
      ) : (
        <div className="applications-list">
          {applications.map((application) => (
            <article
              className="application-card"
              key={application.id}
            >
              <div className="application-card-header">
                <div>
                  <h2>{application.job_title}</h2>
                  <p className="application-company">
                    {application.company}
                  </p>
                </div>

                <StatusBadge status={application.status} />
              </div>

              <div className="application-details">
                {application.date_applied && (
                  <p>
                    <strong>Date applied:</strong>{" "}
                    {application.date_applied}
                  </p>
                )}

                {application.deadline && (
                  <p>
                    <strong>Deadline:</strong>{" "}
                    {application.deadline}
                  </p>
                )}

                {application.job_url && (
                  <p>
                    <strong>Job link:</strong>{" "}
                    <a
                      href={application.job_url}
                      target="_blank"
                      rel="noreferrer"
                    >
                      View job posting
                    </a>
                  </p>
                )}

                {application.notes && (
                  <p className="application-notes">
                    {application.notes}
                  </p>
                )}
              </div>

              <div className="application-card-actions">
                <label>
                  <span>Update status</span>
                  <select
                    value={application.status}
                    onChange={(event) =>
                      void handleStatusChange(
                        application,
                        event.target.value as ApplicationStatus
                      )
                    }
                    disabled={deletingId === application.id}
                  >
                    <option value="Saved">Saved</option>
                    <option value="Applied">Applied</option>
                    <option value="Interview">Interview</option>
                    <option value="Offer">Offer</option>
                    <option value="Rejected">Rejected</option>
                  </select>
                </label>

                <div className="application-action-buttons">
                  <Button
                    variant="secondary"
                    onClick={() => handleEdit(application)}
                  >
                    Edit
                  </Button>

                  <Button
                    variant="quiet"
                    disabled={deletingId === application.id}
                    onClick={() => void handleDelete(application)}
                  >
                    {deletingId === application.id
                      ? "Deleting..."
                      : "Delete"}
                  </Button>
                </div>
              </div>
            </article>
          ))}
        </div>
      )}
    </section>
  );
}
 