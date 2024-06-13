<template>
    <SideBar></SideBar>
    <div class="p-4 sm:ml-64">
        <div class="mx-4 my-4 flex items-center justify-end md:mx-6">
            <Popover v-slot="{ open }" class="relative">
                <PopoverButton
                    :class="open ? 'text-white' : 'text-white/90'"
                    class="group inline-flex items-center rounded-md px-3 py-2 text-base font-medium hover:text-white focus:outline-none focus-visible:ring-2 focus-visible:ring-white/75"
                >
                    <BellIcon
                        :class="open ? 'text-orange-300' : 'text-orange-300/70'"
                        class="ml-2 h-5 w-5 transition duration-150 ease-in-out group-hover:text-blue-400"
                        aria-hidden="true"
                    />
                </PopoverButton>

                <transition
                    enter-active-class="transition duration-200 ease-out"
                    enter-from-class="translate-y-1 opacity-0"
                    enter-to-class="translate-y-0 opacity-100"
                    leave-active-class="transition duration-150 ease-in"
                    leave-from-class="translate-y-0 opacity-100"
                    leave-to-class="translate-y-1 opacity-0"
                >
                    <PopoverPanel
                        class="absolute right-0 z-10 mt-3 origin-top-right w-screen max-w-sm transform px-4 sm:px-0 lg:max-w-3xl"
                    >
                        <div
                            class="overflow-hidden rounded-lg shadow-lg ring-1 ring-black/5"
                        >
                            <div
                                class="relative grid gap-8 bg-white p-7 lg:grid-cols-2"
                            >
                                <a
                                    v-for="item in solutions"
                                    :key="item.name"
                                    :href="item.href"
                                    class="-m-3 flex items-center rounded-lg p-2 transition duration-150 ease-in-out hover:bg-gray-50 focus:outline-none focus-visible:ring focus-visible:ring-orange-500/50"
                                >
                                    <div
                                        class="flex h-10 w-10 shrink-0 items-center justify-center text-white sm:h-12 sm:w-12"
                                    >
                                        <div v-html="item.icon"></div>
                                    </div>
                                    <div class="ml-4">
                                        <p
                                            class="text-sm font-medium text-gray-900"
                                        >
                                            {{ item.name }}
                                        </p>
                                        <p class="text-sm text-gray-500">
                                            {{ item.description }}
                                        </p>
                                    </div>
                                </a>
                            </div>
                            <div class="bg-gray-50 p-4">
                                <a
                                    href="##"
                                    class="flow-root rounded-md px-2 py-2 transition duration-150 ease-in-out hover:bg-gray-100 focus:outline-none focus-visible:ring focus-visible:ring-orange-500/50"
                                >
                                    <span class="flex items-center">
                                        <span
                                            class="text-sm font-medium text-gray-900"
                                        >
                                            Documentation
                                        </span>
                                    </span>
                                    <span class="block text-sm text-gray-500">
                                        Start integrating products and tools
                                    </span>
                                </a>
                            </div>
                        </div>
                    </PopoverPanel>
                </transition>
            </Popover>

            <!-- Profile dropdown -->
            <Menu as="div" class="relative ml-3">
                <div>
                    <MenuButton
                        class="relative flex max-w-xs items-center rounded-full bg-blue-400 text-sm focus:outline-none focus:ring-2 focus:ring-white focus:ring-offset-2 focus:ring-offset-blue-600"
                    >
                        <span class="absolute -inset-1.5" />
                        <span class="sr-only">Open user menu</span>
                        <img
                            class="h-8 w-8 rounded-full"
                            :src="user?.avatar ? user.avatar : nullAvatar"
                            alt="User avatar"
                        />
                    </MenuButton>
                </div>
                <transition
                    enter-active-class="transition ease-out duration-100"
                    enter-from-class="transform opacity-0 scale-95"
                    enter-to-class="transform opacity-100 scale-100"
                    leave-active-class="transition ease-in duration-75"
                    leave-from-class="transform opacity-100 scale-100"
                    leave-to-class="transform opacity-0 scale-95"
                >
                    <MenuItems
                        class="absolute right-0 z-10 mt-2 w-48 origin-top-right rounded-md bg-white py-1 shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none"
                    >
                        <MenuItem
                            v-for="item in userNavigation"
                            :key="item.name"
                            v-slot="{ active }"
                        >
                            <span
                                @click="item.action"
                                :class="[
                                    active ? 'bg-gray-100' : '',
                                    'block px-4 py-2 text-sm text-gray-700',
                                ]"
                                >{{ item.name }}</span
                            >
                        </MenuItem>
                    </MenuItems>
                </transition>
            </Menu>
        </div>
        <slot>
            <RouterView></RouterView>
        </slot>
    </div>
</template>

<script setup lang="ts">
import {
    Menu,
    MenuButton,
    MenuItem,
    MenuItems,
    Popover,
    PopoverButton,
    PopoverPanel,
} from "@headlessui/vue";
import { BellIcon } from "@heroicons/vue/24/outline";
import SideBar from "./components/SideBar.vue";
import { useRouter } from "vue-router";
import { useUserStore } from "@/stores/userStore";
import nullAvatar from "@/assets/vecteezy_illustration-of-human-icon-vector-user-symbol-icon-modern_8442086.jpg";
import { useServerStore } from "@/stores/serverStore";

const userStore = useUserStore();

const router = useRouter();

const userNavigation = [
    { name: "Your Profile", action: () => {} },
    { name: "Settings", action: () => {} },
    {
        name: "Sign out",
        action: () => {
            sessionStorage.clear();
            useUserStore().clear();
            useServerStore().clear();
            router.push({ name: "login" });
        },
    },
];
const solutions: any[] = [];
const user = userStore.currentUser;
</script>
