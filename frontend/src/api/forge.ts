import { api } from './client';
import type { Repo } from '../types/Repo';

export const forgeApi = {
  getRepos: () =>
    api.get<Repo[]>('/forge/repos'),
};
