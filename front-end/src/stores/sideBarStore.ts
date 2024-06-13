import type { SidebarItem } from "@/components/layouts/components/interfaces";
import { defineStore } from "pinia";
import { computed, ref } from "vue";

export const useSideBarStore = defineStore("sidebar", () => {
    const sidebar = ref<SidebarItem[]>([]);

    const checkSideBar = computed(() => {
        return sidebar;
    });

    const selectedSidebar = ref("dashboard");
    const selectedSidebarComputed = computed(() => {
        return selectedSidebar;
    });

    function updateSideBar(value: SidebarItem[]) {
        sidebar.value = value;
    }

    function updateSelectedSidebar(value: string) {
        selectedSidebar.value = value;
    }

    return {
        sidebar,
        updateSideBar,
        checkSideBar,
        selectedSidebarComputed,
        updateSelectedSidebar,
    };
});
