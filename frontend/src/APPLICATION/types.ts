export type ApplicationStatus =
  | "Saved"
  | "Applied"
  | "Interview"
  | "Offer"
  | "Rejected";

export interface Application {
  id: number;
  company: string;
  job_title: string;
  job_url: string | null;
  date_applied: string | null;
  deadline: string | null;
  status: ApplicationStatus;
  notes: string | null;
  created_at: string;
  updated_at: string;
}