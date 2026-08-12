export interface App {
  id: number;
  name: string;
  description: string;
  kind: string;
  language: string;
  repoUrl: string;
  ciUrl: string;
  projectId: number;
  projectName?: string;
  templateId?: number;
  templateName?: string;
  creatorId: number;
  creatorUsername: string;
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

export interface DocPage {
  path: string;
  html: string;
}

export interface TOCItem {
  title: string;
  anchor: string;
  level: number;
}

export interface NavNode {
  title: string;
  path?: string;
  // null for a plain leaf page, an array (possibly empty) for a directory —
  // "is this a folder" is independent of whether it has sibling files.
  children: NavNode[] | null;
}

export interface RenderedDocs {
  pages: DocPage[];
  nav: NavNode[];
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
