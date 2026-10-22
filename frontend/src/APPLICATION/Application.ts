
import { API_BASE_URL } from "../config";
import type { Application } from "./types";

const API_URL = `${API_BASE_URL}/api/jobs`;

export async function getApplications(): Promise<Application[]> {
  const response = await fetch(API_URL);

  if (!response.ok) {
    throw new Error("Failed to load applications.");
  }

  const result = await response.json();

  return result.data ?? result;
}

export async function createApplication(
  application: Omit<Application, "id" | "created_at" | "updated_at">
): Promise<Application> {
  const response = await fetch(API_URL, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(application),
  });

  if (!response.ok) {
    throw new Error("Failed to create application.");
  }

  const result = await response.json();

  return result.data ?? result;
}

export async function updateApplication(
  id: number,
  updates: Partial<Omit<Application, "id" | "created_at" | "updated_at">>
): Promise<Application> {
  const response = await fetch(`${API_URL}/${id}`, {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(updates),
  });

  if (!response.ok) {
    throw new Error("Failed to update application.");
  }

  const result = await response.json();

  return result.data ?? result;
}

export async function deleteApplication(id: number): Promise<void> {
  const response = await fetch(`${API_URL}/${id}`, {
    method: "DELETE",
  });

  if (!response.ok) {
    throw new Error("Failed to delete application.");
  }
}

export interface DashboardStats {
  total_applications: number;
  status_counts: {
    Saved: number;
    Applied: number;
    Interview: number;
    Offer: number;
    Rejected: number;
  };
  upcoming_deadlines: number;
  overdue_applications: number;
  total_interviews: number;
  total_offers: number;
}

export async function getDashboardStats(): Promise<DashboardStats> {
  const response = await fetch(`${API_BASE_URL}/api/dashboard`);

  if (!response.ok) {
    throw new Error("Failed to load dashboard statistics.");
  }

  const result = await response.json();

  return result.data ?? result;
}
