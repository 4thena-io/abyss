import { createRouter, createWebHistory } from 'vue-router';
import Home from './views/Home.vue';
import Apps from './views/Apps.vue';
import Projects from './views/Projects.vue';
import Templates from './views/Templates.vue';
import TemplateDetail from './views/TemplateDetail.vue';
import AppDetail from './views/AppDetail.vue';
import ProjectDetail from './views/ProjectDetail.vue';
import Teams from './views/Teams.vue';
import TeamDetail from './views/TeamDetail.vue';
import Setup from './views/Setup.vue';
import Login from './views/Login.vue';
import Settings from './views/Settings.vue';
import UserProfile from './views/UserProfile.vue';
import ServiceUnavailable from './views/ServiceUnavailable.vue';
import { checkAuthenticated, checkConfigured } from './api/session';

const routes = [
  { path: '/unavailable', component: ServiceUnavailable },
  { path: '/setup', component: Setup },
  { path: '/login', component: Login },
  { path: '/', component: Home },
  { path: '/teams', component: Teams },
  { path: '/teams/:id', component: TeamDetail },
  { path: '/apps', component: Apps },
  { path: '/apps/:id', component: AppDetail },
  { path: '/projects', component: Projects },
  { path: '/projects/:id', component: ProjectDetail },
  { path: '/templates', component: Templates },
  { path: '/templates/:id', component: TemplateDetail },
  { path: '/settings', component: Settings },
  { path: '/users/:id', component: UserProfile },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

// Caches reset on every full page load (which is how login/logout work).
let setupChecked: boolean | null = null;
let authChecked: boolean | null = null;

async function isConfigured(): Promise<boolean> {
  if (setupChecked !== null) return setupChecked;
  setupChecked = await checkConfigured();
  return setupChecked;
}

async function isAuthenticated(): Promise<boolean> {
  if (authChecked !== null) return authChecked;
  authChecked = await checkAuthenticated();
  return authChecked;
}

router.beforeEach(async (to) => {
  if (to.path === '/unavailable') return true;
  try {
    const configured = await isConfigured();

    if (to.path === '/setup') {
      return configured ? '/login' : true;
    }

    if (!configured) return '/setup';

    if (to.path === '/login') {
      const authenticated = await isAuthenticated();
      return authenticated ? '/' : true;
    }

    const authenticated = await isAuthenticated();
    if (!authenticated) return '/login';

    return true;
  } catch {
    return '/unavailable';
  }
});

export default router;
