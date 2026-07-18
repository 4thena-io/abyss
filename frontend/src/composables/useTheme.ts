import { ref, watch } from 'vue';

type Theme = 'light' | 'dark';

const stored = localStorage.getItem('theme') as Theme | null;
const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
const theme = ref<Theme>(stored ?? (prefersDark ? 'dark' : 'light'));

watch(
  theme,
  (val) => {
    document.documentElement.setAttribute('data-theme', val);
    localStorage.setItem('theme', val);
  },
  { immediate: true },
);

export function useTheme() {
  function toggle() {
    theme.value = theme.value === 'light' ? 'dark' : 'light';
  }
  return { theme, toggle };
}
