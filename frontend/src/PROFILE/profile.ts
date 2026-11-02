import { API_BASE_URL } from "../config";

export interface Profile {
  id: string;
  full_name: string | null;
  email: string | null;
  phone: string | null;
  linkedin_url: string | null;
  portfolio_url: string | null;
  cv_text: string | null;
  cv_path: string | null;
  created_at: string;
  updated_at: string;
}

export async function getProfile(): Promise<Profile | null> {
  const res = await fetch(`${API_BASE_URL}/api/profile`);
  if (!res.ok) throw new Error("Failed to load profile");
  const r = await res.json();
  return r.data ?? null;
}

export async function updateProfile(p: Partial<Profile>): Promise<Profile> {
  const res = await fetch(`${API_BASE_URL}/api/profile`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(p),
  });
  if (!res.ok) throw new Error("Failed to update profile");
  const r = await res.json();
  return r.data;
}
