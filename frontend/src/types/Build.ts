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
