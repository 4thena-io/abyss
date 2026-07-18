export interface Project {
  id: number;
  name: string;
  description: string;
  teamId?: number;
  teamName?: string;
  creatorId: number;
  creatorUsername: string;
}
