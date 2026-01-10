import { createRouter, createWebHistory } from "vue-router";
import Home from './views/Home.vue'
import Apps from "./views/Apps.vue";
import Projects from "./views/Projects.vue";
import Templates from "./views/Templates.vue";


const routes = [
	{ path: "/", component: Home },
  { path: "/apps", component: Apps },
  { path: "/projects", component: Projects },
  { path: "/templates", component: Templates }
];

const router = createRouter({
	history: createWebHistory(),
	routes,
});

export default router;
