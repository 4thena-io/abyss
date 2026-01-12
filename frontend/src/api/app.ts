import { api } from './client';
import type { App } from '../types/App';
import type { Build } from '../types/Build';

export interface CreateAppRequest {
  name: string;
  description: string;
  kind: string;
  language: string;
  projectId: number;
  templateId: number;
}

export const appsApi = {
  getAll: () => api.get<App[]>('/apps'),
  getById: (id: number) => api.get<App>(`/apps/${id}`),
  create: (data: CreateAppRequest) => api.post<App>('/apps', data),
  delete: (id: number) => api.delete(`/apps/${id}`),
  getBuilds: (id: number) => api.get<Build[]>(`/apps/${id}/builds`),
};
