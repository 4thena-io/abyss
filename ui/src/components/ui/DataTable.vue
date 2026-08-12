<template>
  <div class="bg-panel border border-b1 rounded-lg overflow-hidden">
    <table class="w-full">
      <thead>
        <tr class="text-left text-t3 text-sm border-b border-b1">
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
      <tbody class="divide-y divide-b1">
        <tr
          v-for="(row, index) in rows"
          :key="rowKey ? row[rowKey] : index"
          class="hover:bg-hover transition-colors"
        >
          <td
            v-for="column in columns"
            :key="column.key"
            class="px-5 py-4 text-t2"
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
  // eslint-disable-next-line @typescript-eslint/no-explicit-any -- generic table accepts heterogeneous row shapes; follow-up: type per call site or make DataTable generic over T
  rows: Record<string, any>[];
  rowKey?: string;
}>();
</script>
