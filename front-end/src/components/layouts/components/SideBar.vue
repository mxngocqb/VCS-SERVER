<template>
    <div class="block">
        <div class="w-64 flex mx-4 my-4 fixed sm:translate-x-0">
            <a href="#" class="flex items-center space-x-3 rtl:space-x-reverse">
                <img :src="logo" class="h-8" alt="VCServer Logo" />
                <span
                    class="self-center text-2xl text-blue-500 font-bold whitespace-nowrap dark:text-white"
                    >ServerMS</span
                >
            </a>
            <button
                data-drawer-target="sidebar-multi-level-sidebar"
                data-drawer-toggle="sidebar-multi-level-sidebar"
                aria-controls="sidebar-multi-level-sidebar"
                type="button"
                class="inline-flex items-center p-2 mt-2 ms-3 text-sm text-gray-500 rounded-lg sm:hidden hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-gray-200 dark:text-gray-400 dark:hover:bg-gray-700 dark:focus:ring-gray-600"
            >
                <span class="sr-only">Open sidebar</span>
                <svg
                    class="w-6 h-6"
                    aria-hidden="true"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                    xmlns="http://www.w3.org/2000/svg"
                >
                    <path
                        clip-rule="evenodd"
                        fill-rule="evenodd"
                        d="M2 4.75A.75.75 0 012.75 4h14.5a.75.75 0 010 1.5H2.75A.75.75 0 012 4.75zm0 10.5a.75.75 0 01.75-.75h7.5a.75.75 0 010 1.5h-7.5a.75.75 0 01-.75-.75zM2 10a.75.75 0 01.75-.75h14.5a.75.75 0 010 1.5H2.75A.75.75 0 012 10z"
                    ></path>
                </svg>
            </button>
        </div>
        <aside
            id="sidebar-multi-level-sidebar"
            class="fixed border-r-2 z-40 w-64 h-full transition-transform -translate-x-full translate-y-16 sm:translate-x-0"
            aria-label="Sidebar"
        >
            <div class="h-full px-3 py-4 overflow-y-auto dark:bg-gray-800">
                <ul class="space-y-2">
                    <li
                        v-for="i in sidebar"
                        :key="i.name"
                        @click="updateSelectedSidebarItem(i.name || '')"
                        class="ring-2 ring-blue-100 rounded-xl"
                        :class="
                            selectedSideBarItem === i.name ? 'bg-blue-100' : ''
                        "
                    >
                        <SideBarItem v-if="!i?.child" :item="i"></SideBarItem>
                        <SideBarGroupItem v-else :group="i"></SideBarGroupItem>
                    </li>
                </ul>
            </div>
        </aside>
    </div>
</template>

<script setup lang="ts">
import logo from "@/assets/_995ed77b-bb47-45ba-aa62-42c8bc67e68a-removebg-preview.png";
import { useSideBarStore } from "@/stores/sideBarStore";
import SideBarItem from "./SideBarItem.vue";
import SideBarGroupItem from "./SideBarGroupItem.vue";

const sidebarStore = useSideBarStore();

const sidebar = sidebarStore.checkSideBar;

const selectedSideBarItem = sidebarStore.selectedSidebarComputed;
const updateSelectedSidebarItem = (value: string) => {
    sidebarStore.updateSelectedSidebar(value);
};
</script>
