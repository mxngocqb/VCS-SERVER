import DashboardSideBar from "@/components/layouts/DashboardSideBar.vue";
import { HomeSideBar } from "@/components/layouts/components/constants";
import type { SidebarItem } from "@/components/layouts/components/interfaces";
import LoginPage from "@/components/pages/LoginPage.vue";
import Forbiden from "@/components/pages/Forbiden.vue";
import FileView from "@/components/views/fileview/FileView.vue";
import AdminView from "@/components/views/adminview/AdminView.vue";
import DashboardView from "@/components/views/dashboard/DashboardView.vue";
import ServerView from "@/components/views/serverView/ServerView.vue";
import { SessionStorageKey } from "@/stores/constants";
import { useSideBarStore } from "@/stores/sideBarStore";
import {
    createRouter,
    createWebHistory,
    type RouteRecordRaw,
} from "vue-router";

const routes: RouteRecordRaw[] = [
    {
        path: "/",
        redirect: "/home",
    },
    {
        name: "home",
        path: "/home",
        component: DashboardSideBar,
        redirect: { name: "dashboard" },
        children: [
            { name: "server", path: "server", component: ServerView },
            { name: "file", path: "storage", component: FileView },
            { name: "admin", path: "admin", component: AdminView },
            { name: "dashboard", path: "dashboard", component: DashboardView },
        ],
        meta: {
            sideBar: HomeSideBar,
            requiresAuth: true,
        },
    },
    { name: "login", path: "/login", component: LoginPage },
    { name: "forbiden", path: "/error", component: Forbiden },
];

const router = createRouter({
    history: createWebHistory("/"),
    routes,
});

function isLogin(): boolean {
    if (sessionStorage.getItem(SessionStorageKey.AUTH_TOKEN)) {
        return true;
    }
    return false;
}

router.beforeEach((to, from) => {
    // instead of having to check every route record with
    // to.matched.some(record => record.meta.requiresAuth)
    if (to.meta?.requiresAuth && !isLogin()) {
        // this route requires auth, check if logged in
        // if not, redirect to login page.
        return {
            name: "login",
        };
    }

    if (to.name === "login" && isLogin()) {
        return { name: "home" };
    }

    if (to.meta?.sideBar) {
        const store = useSideBarStore();
        store.updateSideBar(to.meta.sideBar as SidebarItem[]);
    }
});

export default router;
