<template>
  <div class="h-full overflow-hidden">

    <div v-if="loading" class="h-full flex items-center justify-center">
      <Spinner size="lg" />
    </div>

    <EmptyState
      v-else-if="!docs"
      :icon="DocumentTextIcon"
      title="No docs yet"
      message="Push changes to your docs/ folder to trigger the first render."
    />

    <div v-else class="grid h-full border border-b1 rounded-xl overflow-hidden" style="grid-template-columns: 220px 1fr 200px">

      <!-- Left: page list -->
      <nav class="border-r border-b1 bg-panel overflow-y-auto py-5">
        <template v-for="[section, pages] in groupedPages" :key="section">
          <div
            v-if="section !== '_root'"
            class="px-4 pt-3 pb-1 text-[10px] font-semibold tracking-widest uppercase text-t3"
          >
            {{ formatSection(section) }}
          </div>
          <button
            v-for="page in pages"
            :key="page.path"
            @click="activePath = page.path"
            :class="[
              'w-full text-left px-4 py-[7px] text-[13px] border-l-2 transition-colors',
              activePath === page.path
                ? 'text-acc border-acc bg-acc-s font-medium'
                : 'text-t2 border-transparent hover:bg-surface hover:text-t1',
            ]"
          >
            {{ formatName(page.path) }}
          </button>
        </template>
      </nav>

      <!-- Center: rendered content -->
      <div class="overflow-y-auto">
        <div class="max-w-2xl mx-auto px-10 py-10">
          <div v-if="activePage" class="doc-content" v-html="activePage.html" />
        </div>
      </div>

      <!-- Right: on-page TOC -->
      <div class="border-l border-b1 overflow-y-auto px-5 py-7">
        <template v-if="pageToc.length > 0">
          <div class="text-[10px] font-semibold tracking-widest uppercase text-t3 mb-3">
            On this page
          </div>
          <a
            v-for="item in pageToc"
            :key="item.anchor + item.title"
            :href="item.anchor ? `#${item.anchor}` : undefined"
            :class="[
              'block py-1 border-l-2 border-transparent text-t3 hover:text-acc transition-colors',
              item.level === 1 ? 'pl-2 text-[12px]' : 'pl-5 text-[11px]',
            ]"
          >
            {{ item.title }}
          </a>
        </template>
      </div>

    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { DocumentTextIcon } from '@heroicons/vue/24/outline';
import type { RenderedDocs, DocPage, TOCItem } from '../types';
import Spinner from '../../../components/ui/Spinner.vue';
import EmptyState from '../../../components/ui/EmptyState.vue';

const props = defineProps<{ appId: number }>();

const docs = ref<RenderedDocs | null>(null);
const loading = ref(true);
const activePath = ref('');

const groupedPages = computed((): [string, DocPage[]][] => {
  if (!docs.value?.pages) return [];

  const groups = new Map<string, DocPage[]>();
  for (const page of docs.value.pages) {
    const dir = page.path.includes('/') ? page.path.split('/')[0]! : '_root';
    if (!groups.has(dir)) groups.set(dir, []);
    groups.get(dir)!.push(page);
  }

  return [...groups.entries()]
    .sort(([a], [b]) => {
      if (a === '_root') return -1;
      if (b === '_root') return 1;
      // preserve the order sections first appear in the pages array
      const aIdx = docs.value!.pages.findIndex(p => p.path.startsWith(a + '/'));
      const bIdx = docs.value!.pages.findIndex(p => p.path.startsWith(b + '/'));
      return aIdx - bIdx;
    })
    .map(([section, pages]) => [
      section,
      [...pages].sort((a, b) => {
        const aIsIndex = a.path.endsWith('index.md');
        const bIsIndex = b.path.endsWith('index.md');
        if (aIsIndex && !bIsIndex) return -1;
        if (!aIsIndex && bIsIndex) return 1;
        return a.path.localeCompare(b.path);
      }),
    ] as [string, DocPage[]]);
});

const activePage = computed((): DocPage | undefined =>
  docs.value?.pages.find(p => p.path === activePath.value),
);

const pageToc = computed((): TOCItem[] => {
  if (!activePage.value) return [];
  const el = document.createElement('div');
  el.innerHTML = activePage.value.html;
  return Array.from(el.querySelectorAll('h1, h2, h3')).map(h => ({
    title: h.textContent?.trim() ?? '',
    anchor: (h as HTMLElement).id,
    level: parseInt(h.tagName[1]!),
  }));
});

function formatName(path: string): string {
  const file = path.split('/').pop()?.replace(/\.md$/, '') ?? '';
  if (file === 'index') return 'Overview';
  return file.replace(/[-_]/g, ' ').replace(/\b\w/g, c => c.toUpperCase());
}

function formatSection(dir: string): string {
  return dir.replace(/[-_]/g, ' ').replace(/\b\w/g, c => c.toUpperCase());
}

async function fetchDocs() {
  try {
    loading.value = true;
    const res = await fetch(`/api/apps/${props.appId}/docs`);
    if (res.status === 404) {
      docs.value = null;
      return;
    }
    if (!res.ok) throw new Error('Failed to fetch docs');
    docs.value = await res.json();
    const pages = docs.value?.pages ?? [];
    activePath.value = (pages.find(p => p.path === 'index.md') ?? pages[0])?.path ?? '';
  } catch (err) {
    console.error('Failed to fetch docs:', err);
    docs.value = null;
  } finally {
    loading.value = false;
  }
}

onMounted(fetchDocs);
</script>

<style scoped>
.doc-content :deep(h1) {
  font-size: 26px;
  font-weight: 700;
  letter-spacing: -0.5px;
  color: var(--t1);
  margin-bottom: 8px;
  padding-bottom: 14px;
  border-bottom: 1px solid var(--b1);
}
.doc-content :deep(h2) {
  font-size: 17px;
  font-weight: 600;
  color: var(--t1);
  margin: 28px 0 10px;
}
.doc-content :deep(h3) {
  font-size: 14px;
  font-weight: 600;
  color: var(--t1);
  margin: 20px 0 8px;
}
.doc-content :deep(p) {
  font-size: 14px;
  color: var(--t2);
  line-height: 1.75;
  margin-bottom: 14px;
}
.doc-content :deep(a) {
  color: var(--acc);
}
.doc-content :deep(a:hover) {
  text-decoration: underline;
}
.doc-content :deep(code) {
  font-family: var(--mono);
  font-size: 12px;
  background: var(--surface);
  border: 1px solid var(--b1);
  border-radius: 4px;
  padding: 1px 5px;
  color: var(--t1);
}
.doc-content :deep(pre) {
  background: var(--surface);
  border: 1px solid var(--b1);
  border-radius: 8px;
  padding: 16px;
  margin: 14px 0;
  overflow-x: auto;
  font-family: var(--mono);
  font-size: 12px;
  color: var(--t2);
  line-height: 1.65;
}
.doc-content :deep(pre code) {
  background: none;
  border: none;
  padding: 0;
  font-size: inherit;
}
.doc-content :deep(ul),
.doc-content :deep(ol) {
  padding-left: 18px;
  margin-bottom: 14px;
}
.doc-content :deep(li) {
  font-size: 14px;
  color: var(--t2);
  line-height: 1.75;
  margin-bottom: 4px;
}
.doc-content :deep(blockquote) {
  border-left: 3px solid var(--b2);
  padding-left: 16px;
  margin: 14px 0;
  color: var(--t2);
}
.doc-content :deep(hr) {
  border: none;
  border-top: 1px solid var(--b1);
  margin: 24px 0;
}
.doc-content :deep(table) {
  width: 100%;
  border-collapse: collapse;
  margin-bottom: 14px;
  font-size: 13px;
}
.doc-content :deep(th) {
  text-align: left;
  padding: 8px 12px;
  background: var(--surface);
  border: 1px solid var(--b1);
  color: var(--t1);
  font-weight: 600;
}
.doc-content :deep(td) {
  padding: 8px 12px;
  border: 1px solid var(--b1);
  color: var(--t2);
}
</style>
