export interface Deployment {
  id: number;
  environment: 'production' | 'staging' | 'dev';
  status: 'success' | 'failure' | 'running' | 'pending';
  commit: string;
  triggeredBy: string;
  duration: number;
  deployedAt: string;
}
