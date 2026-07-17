import { api } from './client';
import type { Template } from '../features/templates/types';
import type { App } from '../features/apps/types';

export interface CreateTemplateRequest {
  source: 'repo' | 'blank';
  name: string;
  description: string;
  kind: string;
  language: string;
  repoUrl?: string;
}

export interface UpdateTemplateRequest {
  name: string;
  description: string;
  kind: string;
  language: string;
}

export const templatesApi = {
  getAll: () => api.get<Template[]>('/templates'),
  getById: (id: number) => api.get<Template>(`/templates/${id}`),
  getApps: (id: number) => api.get<App[]>(`/templates/${id}/apps`),
  create: (data: CreateTemplateRequest) => api.post<Template>('/templates', data),
  update: (id: number, data: UpdateTemplateRequest) => api.patch<Template>(`/templates/${id}`, data),
  delete: (id: number) => api.delete(`/templates/${id}`),
};
