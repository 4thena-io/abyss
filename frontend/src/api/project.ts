import { api } from './client';
import type { Project } from '../features/projects/types';
import type { App } from '../features/apps/types';

export interface CreateProjectRequest {
  name: string;
  description: string;
  teamId?: number | null;
}

export interface UpdateProjectRequest {
  name: string;
  description: string;
  teamId?: number | null;
}

export const projectsApi = {
  getAll: () => api.get<Project[]>('/projects'),
  getById: (id: number) => api.get<Project>(`/projects/${id}`),
  create: (data: CreateProjectRequest) => api.post<Project>('/projects', data),
  update: (id: number, data: UpdateProjectRequest) => api.patch<Project>(`/projects/${id}`, data),
  delete: (id: number) => api.delete(`/projects/${id}`),
  getApps: (id: number) => api.get<App[]>(`/projects/${id}/apps`),
};
