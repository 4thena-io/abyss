export interface Team {
  id: number;
  name: string;
  description: string;
  memberCount: number;
  projectCount: number;
  appCount: number;
}

export interface TeamMember {
  username: string;
  role: 'owner' | 'member';
  joinedAt: string;
}
