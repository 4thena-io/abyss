import { globalIgnores } from 'eslint/config';
import pluginVue from 'eslint-plugin-vue';
import { defineConfigWithVueTs, vueTsConfigs } from '@vue/eslint-config-typescript';
import eslintConfigPrettier from 'eslint-config-prettier';

export default defineConfigWithVueTs(
  {
    name: 'app/files-to-lint',
    files: ['**/*.{ts,mts,tsx,vue}'],
  },
  globalIgnores([
    '**/dist/**',
    '**/dist-ssr/**',
    '**/coverage/**',
    '**/node_modules/**',
    // Vite's scaffolded ambient module declaration; not hand-edited.
    'src/vite-env.d.ts',
  ]),
  pluginVue.configs['flat/essential'],
  vueTsConfigs.recommended,
  {
    name: 'app/rule-overrides',
    files: ['**/*.vue'],
    rules: {
      // Route views and top-level components (Home.vue, Button.vue, ...) are
      // never registered as custom elements, so there's no collision risk
      // this rule guards against; enforcing it would mean renaming every
      // view/component file for no behavioral benefit.
      'vue/multi-word-component-names': 'off',
    },
  },
  // Must stay last: turns off stylistic ESLint/vue rules that would
  // otherwise fight Prettier's formatting decisions.
  eslintConfigPrettier,
);
