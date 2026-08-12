import { api } from './client';
import type { Team, TeamMember } from '../features/teams/types';
import type { Project } from '../features/projects/types';

export interface CreateTeamRequest {
  name: string;
  description: string;
}

export interface UpdateTeamRequest {
  name: string;
  description: string;
}

export interface AddMemberRequest {
  username: string;
  role: string;
}

export const teamsApi = {
  getAll: () => api.get<Team[]>('/teams'),
  getById: (id: number) => api.get<Team>(`/teams/${id}`),
  create: (data: CreateTeamRequest) => api.post<Team>('/teams', data),
  update: (id: number, data: UpdateTeamRequest) => api.patch<Team>(`/teams/${id}`, data),
  delete: (id: number) => api.delete(`/teams/${id}`),
  getMembers: (id: number) => api.get<TeamMember[]>(`/teams/${id}/members`),
  getProjects: (id: number) => api.get<Project[]>(`/teams/${id}/projects`),
  addMember: (id: number, data: AddMemberRequest) =>
    api.post<TeamMember>(`/teams/${id}/members`, data),
  removeMember: (teamId: number, memberId: number) =>
    api.delete(`/teams/${teamId}/members/${memberId}`),
};
