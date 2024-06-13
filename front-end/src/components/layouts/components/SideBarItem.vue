<template>
    <RouterLink
        :to="{ name: item.name } || '#'"
        class="flex items-center text-base font-medium p-2 text-blue-900 rounded-lg dark:text-white hover:bg-blue-100 dark:hover:bg-blue-700 group"
    >
        <svg
            class="w-6 h-6 text-blue-500 transition duration-75 dark:text-blue-400 group-hover:text-blue-900 dark:group-hover:text-white"
            aria-hidden="true"
            xmlns="http://www.w3.org/2000/svg"
            fill="currentColor"
            viewBox="0 0 22 22"
            v-html="props.item.icon"
        ></svg>
        <span class="flex-1 ms-3 whitespace-nowrap">{{
            props.item.label
        }}</span>
        <span
            v-if="props.item?.extendText"
            class="inline-flex items-center justify-center w-3 h-3 p-3 ms-3 text-sm font-medium text-blue-800 bg-blue-100 rounded-full dark:bg-blue-900 dark:text-blue-300"
            >{{ props.item?.extendText }}</span
        >
    </RouterLink>
</template>
<script setup lang="ts">
import type { SidebarItem } from "./interfaces";

const props = defineProps<{
    item: SidebarItem;
}>();
</script>
