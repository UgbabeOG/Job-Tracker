import { API_BASE_URL } from "../config";

export interface MatchResponse {
  score: number;
  matched_keywords: string[];
  missing_keywords: string[];
  total_keywords: number;
}

export async function matchText(
  cvText: string,
  jobDescription: string,
  requiredKeywords?: string[]
): Promise<MatchResponse> {
  const response = await fetch(`${API_BASE_URL}/api/match`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      cv_text: cvText,
      job_description: jobDescription,
      required_keywords: requiredKeywords,
    }),
  });
  if (!response.ok) {
    const err = await response.json().catch(() => ({}));
    throw new Error(err?.error?.message ?? "Failed to match keywords.");
  }
  const result = await response.json();
  return result.data;
}
