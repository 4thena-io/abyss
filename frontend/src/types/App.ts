export interface App {
  id: string;
  name: string;
  description: string;
  kind: string;
  language: string;
  repo_id: number;
  repo_url: string;
  clone_url: string;
  ci_id: number;
  ci_url: string;
  project: string;
  status: string;
  created_at: string;
  updated_at: string;
}
