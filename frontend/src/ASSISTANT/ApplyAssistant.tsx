import { useState } from "react";
import { parseUrl, getSuggestions } from "./assistant";
import { Button, ErrorMessage } from "../components/ui";

export default function ApplyAssistant() {
  const [url, setUrl] = useState("");
  const [company, setCompany] = useState("");
  const [jobTitle, setJobTitle] = useState("");
  const [save, setSave] = useState(true);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [result, setResult] = useState<any>(null);
  const [sugs, setSugs] = useState<any>(null);

  async function handleParse(e: React.FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError("");
    setResult(null);
    setSugs(null);
    try {
      const data = await parseUrl(url, save, company || undefined, jobTitle || undefined);
      setResult(data);
      if (data.job_id) {
        const s = await getSuggestions(data.job_id);
        setSugs(s);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed");
    } finally {
      setLoading(false);
    }
  }

  function copy(t: string) {
    navigator.clipboard.writeText(t);
  }

  return (
    <section className="applications-page">
      <header className="applications-header">
        <div>
          <h1>Apply Assistant</h1>
          <p>Parse job URLs, get suggestions, and copy answers. You submit manually.</p>
        </div>
      </header>
      {error && <ErrorMessage message={error} />}
      <form onSubmit={handleParse} style={{ display: "grid", gap: "0.75rem", maxWidth: "700px" }}>
        <input placeholder="Paste job URL" value={url} onChange={(e) => setUrl(e.target.value)} required />
        <input placeholder="Company (optional)" value={company} onChange={(e) => setCompany(e.target.value)} />
        <input placeholder="Job title (optional)" value={jobTitle} onChange={(e) => setJobTitle(e.target.value)} />
        <label><input type="checkbox" checked={save} onChange={(e) => setSave(e.target.checked)} /> Save to Applications</label>
        <Button type="submit" disabled={loading}>{loading ? "Parsing..." : "Parse URL"}</Button>
      </form>

      {result && (
        <div style={{ marginTop: "1.5rem" }}>
          <h2>Parsed</h2>
          <p><strong>Platform:</strong> {result.parsed?.source_platform}</p>
          <p><strong>URL:</strong> <a href={result.parsed?.source_url} target="_blank" rel="noreferrer">{result.parsed?.source_url}</a></p>
          {result.job_id && <p><strong>Saved to Applications (ID):</strong> {result.job_id}</p>}
        </div>
      )}

      {sugs && (
        <div style={{ marginTop: "1.5rem" }}>
          <h2>Suggestions</h2>
          <p><strong>Match score:</strong> {sugs.match_score}%</p>
          <p><strong>Matched:</strong> {sugs.matched?.join(", ") || "none"}</p>
          <p><strong>Missing:</strong> {sugs.missing?.join(", ") || "none"}</p>
          <h3>Checklist & Suggestions</h3>
          <ul>
            {(sugs.suggestions || []).map((s: any, i: number) => (
              <li key={i}>
                <strong>{s.field}:</strong> {s.suggestion} {s.needs_review && <em>(review)</em>} <Button variant="quiet" onClick={() => copy(s.suggestion)}>Copy</Button>
              </li>
            ))}
          </ul>
        </div>
      )}
    </section>
  );
}
