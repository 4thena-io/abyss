import { api } from './client';
import type { Project } from '../features/projects/types';
import type { App } from '../features/apps/types';

export interface CreateProjectRequest {
  name: string;
  description: string;
}

export const projectsApi = {
  getAll: () => api.get<Project[]>('/projects'),
  getById: (id: number) => api.get<Project>(`/projects/${id}`),
  create: (data: CreateProjectRequest) => api.post<Project>('/projects', data),
  delete: (id: number) => api.delete(`/projects/${id}`),
  getApps: (id: number) => api.get<App[]>(`/projects/${id}/apps`),
};
