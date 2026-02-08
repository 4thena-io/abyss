<template>
  <div class="bg-gray-800 border border-gray-700 rounded-lg overflow-hidden">
    <table class="w-full">
      <thead>
        <tr class="text-left text-gray-400 text-sm border-b border-gray-700">
          <th 
            v-for="column in columns" 
            :key="column.key" 
            class="px-5 py-3 font-medium"
            :class="column.class"
          >
            {{ column.label }}
          </th>
        </tr>
      </thead>
      <tbody class="divide-y divide-gray-700">
        <tr 
          v-for="(row, index) in rows" 
          :key="rowKey ? row[rowKey] : index" 
          class="hover:bg-gray-750 transition-colors"
        >
          <td 
            v-for="column in columns" 
            :key="column.key" 
            class="px-5 py-4"
            :class="column.cellClass"
          >
            <slot :name="`cell-${column.key}`" :row="row" :value="row[column.key]">
              {{ row[column.key] }}
            </slot>
          </td>
        </tr>
      </tbody>
    </table>

    <slot name="footer" />
  </div>
</template>

<script setup lang="ts">
export interface TableColumn {
  key: string;
  label: string;
  class?: string;
  cellClass?: string;
}

defineProps<{
  columns: TableColumn[];
  rows: Record<string, any>[];
  rowKey?: string;
}>();
</script>
