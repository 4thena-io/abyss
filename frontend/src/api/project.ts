import { api } from './client';
import type { Project } from '../types/Project';
import type { App } from '../types/App';

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
