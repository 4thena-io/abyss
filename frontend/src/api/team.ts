import { api } from './client';
import type { Team, TeamMember } from '../features/teams/types';
import type { Project } from '../features/projects/types';

export interface CreateTeamRequest {
  name: string;
  description: string;
}

export const teamsApi = {
  getAll: () => api.get<Team[]>('/teams'),
  getById: (id: number) => api.get<Team>(`/teams/${id}`),
  create: (data: CreateTeamRequest) => api.post<Team>('/teams', data),
  delete: (id: number) => api.delete(`/teams/${id}`),
  getMembers: (id: number) => api.get<TeamMember[]>(`/teams/${id}/members`),
  getProjects: (id: number) => api.get<Project[]>(`/teams/${id}/projects`),
};
