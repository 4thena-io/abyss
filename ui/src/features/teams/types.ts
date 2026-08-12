export interface Team {
  id: number;
  name: string;
  description: string;
  memberCount: number;
  projectCount: number;
  appCount: number;
}

export interface TeamMember {
  id: number;
  teamId: number;
  userId: number;
  username: string;
  email: string;
  avatarUrl: string;
  role: string;
  joinedAt: string;
}
