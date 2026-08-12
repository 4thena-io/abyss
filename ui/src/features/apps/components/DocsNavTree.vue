<template>
  <template v-for="(node, index) in nodes" :key="node.path ?? node.title">
    <!-- Top-level entries are separated from each other by a rule. -->
    <div v-if="depth === 0 && index > 0" class="border-t border-b1 mx-4 my-2" />

    <!-- Directory: children is an array (possibly empty), independent of
         whether this folder happens to have sibling files — a folder with
         only an index.md looks the same as one with ten siblings. If it
         also has a path, the header is a clickable link to that folder's
         index.md; without a path (no index.md in that directory) it's a
         non-clickable group label. -->
    <template v-if="node.children !== null">
      <button
        v-if="node.path"
        @click="$emit('select', node.path)"
        :style="{ paddingLeft: `${indent}px` }"
        :class="[
          'relative w-full text-left py-[7px] pr-4 text-[13px] transition-colors',
          activePath === node.path ? 'text-acc' : [idleColor, 'hover:text-t1'],
        ]"
      >
        <!-- Depth 0/1 active state is a color change only; depth 2+ also
             highlights the connector line segment next to this row, at the
             parent's indent (where that line already runs), layered over
             the plain grey line drawn below. -->
        <span
          v-if="depth >= 2 && activePath === node.path"
          class="absolute top-1/2 -translate-y-1/2 h-4 w-0.5 bg-acc"
          :style="{ left: `${parentIndent}px` }"
        />
        {{ node.title }}
      </button>
      <div
        v-else
        :style="{ paddingLeft: `${indent}px` }"
        :class="['py-[7px] pr-4 text-[13px]', idleColor]"
      >
        {{ node.title }}
      </div>

      <!-- Nesting below depth 0/1 is shown by indentation + a single
           connector line at the parent's own indent, dropping through its
           visible children — no background panels, one shared background
           throughout. Depth 0 -> 1 stays color-only, no line. -->
      <template v-if="node.children.length">
        <div v-if="depth >= 1" class="relative mt-1">
          <div class="absolute top-0 bottom-0 w-px bg-b2" :style="{ left: `${indent}px` }" />
          <DocsNavTree
            :nodes="node.children"
            :active-path="activePath"
            :depth="depth + 1"
            @select="$emit('select', $event)"
          />
        </div>
        <DocsNavTree
          v-else
          :nodes="node.children"
          :active-path="activePath"
          :depth="depth + 1"
          @select="$emit('select', $event)"
        />
      </template>
    </template>

    <!-- Leaf page -->
    <button
      v-else
      @click="$emit('select', node.path!)"
      :style="{ paddingLeft: `${indent}px` }"
      :class="[
        'relative w-full text-left py-[7px] pr-4 text-[13px] transition-colors',
        activePath === node.path ? 'text-acc' : [idleColor, 'hover:text-t1'],
      ]"
    >
      <span
        v-if="depth >= 2 && activePath === node.path"
        class="absolute top-1/2 -translate-y-1/2 h-4 w-0.5 bg-acc"
        :style="{ left: `${parentIndent}px` }"
      />
      {{ node.title }}
    </button>
  </template>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { NavNode } from '../types';

const props = withDefaults(
  defineProps<{
    nodes: NavNode[];
    activePath: string;
    depth?: number;
  }>(),
  { depth: 0 },
);

defineEmits<{ select: [path: string] }>();

// Depths 0 and 1 share the same indent; depth 2+ steps in by 12px per level.
const indent = computed(() => 16 + Math.max(props.depth - 1, 0) * 12);

// The x-position of the connector line this row hangs from (drawn by the
// parent, at the parent's own indent) — used to align the active-row
// highlight with that existing line instead of the row's own text indent.
const parentIndent = computed(() => 16 + Math.max(props.depth - 2, 0) * 12);

// Top-level entries are bright; everything nested underneath one of them
// (any depth 1+) is a single, uniformly muted grey — not progressively
// darker per level.
const idleColor = computed(() => (props.depth === 0 ? 'text-t1' : 'text-t2'));
</script>
