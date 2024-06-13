import type { IUser } from "@/plugins/axios/user/interfaces";
import { defineStore } from "pinia";
import { computed, ref } from "vue";

export const useUserStore = defineStore("user", () => {
    const user = ref<IUser>();

    const currentUser = computed(() => user);

    const listUser = ref<IUser[]>();

    const listUserComputed = computed(() => {
        return listUser;
    });

    const selectedUser = ref<IUser>();

    const selectedUserComputed = computed(() => {
        return selectedUser;
    });

    function setUser(u: IUser) {
        user.value = u;
    }

    function setListUser(value: IUser[]) {
        listUser.value = value;
    }

    function setSelectedUser(id: number) {
        selectedUser.value = listUser.value?.find((e) => e.id === id);
    }

    function clear() {
        user.value = undefined;
    }

    return {
        currentUser,
        setUser,
        clear,
        listUserComputed,
        setListUser,
        selectedUserComputed,
        setSelectedUser,
    };
});
