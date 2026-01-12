import { api } from './client';
import type { Template } from '../types/Template';

export interface CreateTemplateRequest {
  name: string;
  repoUrl: string;
  kind: string;
  language: string;
  description: string;
}

export const templatesApi = {
  getAll: () => api.get<Template[]>('/templates'),
  getById: (id: number) => api.get<Template>(`/templates/${id}`),
  create: (data: CreateTemplateRequest) => api.post<Template>('/templates', data),
  delete: (id: number) => api.delete(`/templates/${id}`),
};
