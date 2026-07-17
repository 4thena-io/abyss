export interface Project {
  id: number;
  name: string;
  description: string;
  teamId?: number;
  creatorId: number;
  creatorUsername: string;
}
