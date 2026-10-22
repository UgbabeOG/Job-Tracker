
import { useCallback, useEffect, useState } from "react";
import { NavLink } from "react-router-dom";
import {
  StatusBadge,
  LoadingState,
  ErrorMessage,
} from "../components/ui";
import type { ApplicationStatus } from "./types";
import type { DashboardStats } from "./Application";
import { getDashboardStats } from "./Application";
import "./Dashboard.css";

const STATUS_ORDER: ApplicationStatus[] = [
  "Saved",
  "Applied",
  "Interview",
  "Offer",
  "Rejected",
];

type StatCardProps = {
  label: string;
  value: number;
  variant?: "default" | "upcoming" | "overdue";
  indicator?: string;
};

function StatCard({
  label,
  value,
  variant = "default",
  indicator,
}: StatCardProps) {
  return (
    <article className={`stat-card stat-card--${variant}`}>
      <header className="stat-card-header">
        <span className="stat-label">{label}</span>
        {indicator && (
          <span className={`stat-indicator stat-indicator--${variant}`}>
            {indicator}
          </span>
        )}
      </header>
      <p className="stat-value">{value}</p>
    </article>
  );
}

export default function Dashboard() {
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const loadStats = useCallback(async () => {
    setLoading(true);
    setError("");

    try {
      const data = await getDashboardStats();
      setStats(data);
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : "Failed to load dashboard statistics."
      );
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void loadStats();
  }, [loadStats]);

  if (loading) {
    return <LoadingState label="Loading dashboard..." />;
  }

  if (error) {
    return <ErrorMessage message={error} />;
  }

  const statusCounts = stats?.status_counts ?? {};

  return (
    <section className="dashboard-page">
      <header className="dashboard-header">
        <div>
          <h1>Dashboard</h1>
          <p className="dashboard-description">
            Your job-search overview at a glance.
          </p>
        </div>

        <NavLink
          to="/applications"
          className={({ isActive }) =>
            `button button-secondary ${isActive ? "nav-link-active" : ""}`
          }
        >
          View all applications
        </NavLink>
      </header>

      {error && <ErrorMessage message={error} />}

      <div className="stats-grid">
        <StatCard
          label="Total applications"
          value={stats?.total_applications ?? 0}
        />

        <StatCard
          label="Upcoming deadlines"
          value={stats?.upcoming_deadlines ?? 0}
          variant="upcoming"
          indicator="Upcoming"
        />

        <StatCard
          label="Overdue applications"
          value={stats?.overdue_applications ?? 0}
          variant="overdue"
          indicator="Overdue"
        />

        <StatCard
          label="Total interviews"
          value={stats?.total_interviews ?? 0}
        />

        <StatCard
          label="Total offers"
          value={stats?.total_offers ?? 0}
        />
      </div>

      <section className="status-breakdown">
        <h2>Applications by status</h2>
        <div className="status-list">
          {STATUS_ORDER.map((status) => (
            <article
              className="status-item"
              key={status}
            >
              <div className="status-item-main">
                <StatusBadge status={status} />
                <span className="status-name">{status}</span>
              </div>
              <span className="status-count">
                {(statusCounts as Record<string, number>)[status] ?? 0}
              </span>
            </article>
          ))}
        </div>
      </section>
    </section>
  );
}
