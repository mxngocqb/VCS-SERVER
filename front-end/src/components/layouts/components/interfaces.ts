import type { Component } from "vue";

export interface SidebarItem {
    name?: string;
    label: string;
    extendText?: string;
    icon?: string;
    child?: SidebarItem[];
}
