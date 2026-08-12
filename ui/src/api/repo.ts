import { api } from './client';
import type { Repo } from '../types/Repo';

export const repoApi = {
  getRepos: () => api.get<Repo[]>('/repos'),
};
