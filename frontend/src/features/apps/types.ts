export interface App {
  id: number;
  name: string;
  description: string;
  kind: string;
  language: string;
  repoUrl: string;
  ciUrl: string;
  projectId: string;
  templateId: string;
}

export interface Build {
  id: number;
  number: number;
  status: string;
  branch: string;
  commit: string;
  duration: number;
  startedAt: string;
  link: string;
}

export interface Deployment {
  id: number;
  environment: 'production' | 'staging' | 'dev';
  status: 'success' | 'failure' | 'running' | 'pending';
  commit: string;
  triggeredBy: string;
  duration: number;
  deployedAt: string;
}
