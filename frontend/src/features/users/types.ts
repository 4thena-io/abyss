export interface UserProfile {
  id: number;
  username: string;
  email: string;
  avatarUrl: string;
  isAdmin: boolean;
  teamCount: number;
  projectCount: number;
  appCount: number;
  teams: UserTeamMembership[];
}

export interface UserTeamMembership {
  id: number;
  name: string;
  role: string;
  memberCount: number;
  projectCount: number;
}
