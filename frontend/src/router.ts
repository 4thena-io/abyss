import { createRouter, createWebHistory } from "vue-router";
import Home from './views/Home.vue'
import Apps from "./views/Apps.vue";


const routes = [
	{ path: "/", component: Home },
  { path: "/apps", component: Apps }
];

const router = createRouter({
	history: createWebHistory(),
	routes,
});

export default router;
