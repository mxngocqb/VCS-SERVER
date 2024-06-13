<template>
    <button
        type="button"
        class="flex font-medium items-center w-full p-2 text-base text-gray-900 transition duration-75 rounded-lg group hover:bg-gray-100 dark:text-white dark:hover:bg-gray-700"
        aria-controls="dropdown-example"
        data-collapse-toggle="dropdown-example"
    >
        <svg
            class="w-6 h-6 text-gray-500 transition duration-75 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white"
            aria-hidden="true"
            xmlns="http://www.w3.org/2000/svg"
            fill="currentColor"
            viewBox="0 0 22 22"
            v-html="props.group.icon"
        ></svg>
        <span class="flex-1 ms-3 text-left rtl:text-right whitespace-nowrap">{{
            props.group.label
        }}</span>
        <svg
            class="w-3 h-3"
            aria-hidden="true"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 10 6"
        >
            <path
                stroke="currentColor"
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="m1 1 4 4 4-4"
            />
        </svg>
    </button>
    <ul id="dropdown-example" class="hidden py-2 space-y-2">
        <li v-for="item in props.group.child" :key="item.label">
            <RouterLink
                class="flex font-medium items-center w-full p-2 text-gray-900 transition duration-75 rounded-lg pl-11 group hover:bg-gray-100 dark:text-white dark:hover:bg-gray-700"
                :to="item?.name || '#'"
                >{{ item.label }}</RouterLink
            >
        </li>
    </ul>
</template>
<script setup lang="ts">
import type { SidebarItem } from "./interfaces";

const props = defineProps<{
    group: SidebarItem;
}>();
</script>
