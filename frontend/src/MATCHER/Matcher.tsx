import { useEffect, useState } from "react";
import { getApplications } from "../APPLICATION/Application";
import type { Application } from "../APPLICATION/types";
import { matchText, type MatchResponse } from "./match";
import { Button, LoadingState, ErrorMessage, EmptyState } from "../components/ui";

export default function Matcher() {
  const [cvText, setCvText] = useState("");
  const [applications, setApplications] = useState<Application[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [results, setResults] = useState<Record<number, MatchResponse | string>>({});
  const [running, setRunning] = useState(false);

  useEffect(() => {
    getApplications()
      .then(setApplications)
      .catch(() => setError("Failed to load applications."))
      .finally(() => setLoading(false));
  }, []);

  async function runAll() {
    if (!cvText.trim()) {
      setError("CV text is required.");
      return;
    }
    setRunning(true);
    setError("");
    const next: Record<number, MatchResponse | string> = {};
    for (const app of applications) {
      const description = app.job_description?.trim() || "";
      if (!description) {
        next[app.id] = "No job description available for this job.";
        continue;
      }
      try {
        next[app.id] = await matchText(cvText, description);
      } catch (err) {
        next[app.id] = err instanceof Error ? err.message : "Match failed.";
      }
    }
    setResults(next);
    setRunning(false);
  }

  if (loading) return <LoadingState />;
  if (error && applications.length === 0) return <ErrorMessage message={error} />;

  return (
    <section className="applications-page">
      <header className="applications-header">
        <div>
          <h1>Keyword Matcher</h1>
          <p>Compare your CV text against your saved jobs.</p>
        </div>
      </header>
      {error && <ErrorMessage message={error} />}
      <textarea
        value={cvText}
        onChange={(e) => setCvText(e.target.value)}
        placeholder="Paste your CV text here"
        rows={8}
        style={{ width: "100%", marginBottom: "1rem" }}
      />
      <Button onClick={runAll} disabled={running || applications.length === 0}>
        {running ? "Matching..." : "Match all jobs"}
      </Button>

      {applications.length === 0 ? (
        <EmptyState title="No applications" description="Add jobs first." />
      ) : (
        <div className="applications-list" style={{ marginTop: "1rem" }}>
          {applications.map((app) => {
            const r = results[app.id];
            return (
              <article className="application-card" key={app.id}>
                <div className="application-card-header">
                  <div>
                    <h2>{app.job_title}</h2>
                    <p className="application-company">{app.company}</p>
                  </div>
                </div>
                <Button
                  onClick={async () => {
                    if (!cvText.trim()) { setError("CV text is required."); return; }
                    const app2 = app;
                    const description = app2.job_description?.trim() || "";
                    if (!description) { setResults((prev) => ({ ...prev, [app2.id]: "No job description available for this job." })); return; }
                    try {
                      const res = await matchText(cvText, description);
                      setResults((prev) => ({ ...prev, [app2.id]: res }));
                    } catch (err) {
                      setResults((prev) => ({ ...prev, [app2.id]: err instanceof Error ? err.message : "Match failed." }));
                    }
                  }}
                >
                  Match this job
                </Button>
                {typeof r === "string" ? (
                  <p>{r}</p>
                ) : r ? (
                  <div>
                    <p><strong>Score:</strong> {r.score}%</p>
                    <p><strong>Matched:</strong> {r.matched_keywords.join(", ") || "none"}</p>
                    <p><strong>Missing:</strong> {r.missing_keywords.join(", ") || "none"}</p>
                    <p><strong>Total keywords:</strong> {r.total_keywords}</p>
                  </div>
                ) : (
                  <p>Not matched yet.</p>
                )}
              </article>
            );
          })}
        </div>
      )}
    </section>
  );
}
