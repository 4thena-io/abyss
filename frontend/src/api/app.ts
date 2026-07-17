import { api } from './client';
import type { App, Build, Deployment } from '../features/apps/types';

export interface CreateAppRequest {
  name: string;
  description: string;
  kind: string;
  language: string;
  projectId?: number;
  templateId?: number;
  repoId?: number;
}

export interface UpdateAppRequest {
  name: string;
  description: string;
  projectId?: number | null;
}

export const appsApi = {
  getAll: () => api.get<App[]>('/apps'),
  getById: (id: number) => api.get<App>(`/apps/${id}`),
  create: (data: CreateAppRequest) => api.post<App>('/apps', data),
  update: (id: number, data: UpdateAppRequest) => api.patch<App>(`/apps/${id}`, data),
  delete: (id: number) => api.delete(`/apps/${id}`),
  getBuilds: (id: number) => api.get<Build[]>(`/apps/${id}/builds`),
  getDeployments: (id: number) => api.get<Deployment[]>(`/apps/${id}/deployments`),
};
