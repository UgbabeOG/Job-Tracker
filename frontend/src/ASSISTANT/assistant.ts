import { API_BASE_URL } from "../config";

export async function parseUrl(url: string, saveToTracker: boolean, company?: string, jobTitle?: string) {
  const res = await fetch(`${API_BASE_URL}/api/assistant/parse-url`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ url, save_to_tracker: saveToTracker, company, job_title: jobTitle }),
  });
  if (!res.ok) {
    const e = await res.json().catch(() => ({}));
    throw new Error(e?.error?.message || "Failed to parse");
  }
  const r = await res.json();
  return r.data;
}

export async function getSuggestions(jobId: number | string) {
  const res = await fetch(`${API_BASE_URL}/api/assistant/suggestions/${jobId}`);
  if (!res.ok) throw new Error("Failed to get suggestions");
  const r = await res.json();
  return r.data;
}

export async function getAudit() {
  const res = await fetch(`${API_BASE_URL}/api/assistant/audit`);
  if (!res.ok) throw new Error("Failed to get audit");
  const r = await res.json();
  return r.data;
}
