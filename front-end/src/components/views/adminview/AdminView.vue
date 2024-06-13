<template>
    <div class="relative overflow-x-auto sm:rounded-lg">
        <div
            class="flex flex-column sm:flex-row flex-wrap space-y-4 sm:space-y-0 items-center justify-between pb-4"
        >
            <div class="flex">
                <button
                    @click="isOpenCreateUserPopup = !isOpenCreateUserPopup"
                    type="button"
                    class="text-white bg-[#1da1f2] hover:bg-[#1da1f2]/90 focus:ring-4 focus:outline-none focus:ring-[#1da1f2]/50 font-medium rounded-lg text-sm px-5 py-2.5 text-center inline-flex items-center dark:focus:ring-[#1da1f2]/55 mx-1"
                >
                    New User
                </button>
            </div>
        </div>
        <EmptyStates v-if="!listUser?.length"></EmptyStates>
        <table
            v-else
            class="w-full text-sm text-left rtl:text-right text-gray-500 dark:text-gray-400"
        >
            <thead
                class="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400"
            >
                <tr>
                    <th scope="col" class="p-4">
                        <div class="flex items-center">
                            <input
                                id="checkbox-all-search"
                                type="checkbox"
                                class="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 dark:focus:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600"
                            />
                            <label for="checkbox-all-search" class="sr-only"
                                >checkbox</label
                            >
                        </div>
                    </th>
                    <th scope="col" class="px-6 py-3">Full Name</th>
                    <th scope="col" class="px-6 py-3">Create By</th>
                    <th scope="col" class="px-6 py-3">Create At</th>
                    <th scope="col" class="px-6 py-3">Email</th>
                    <th scope="col" class="px-6 py-3">Phone</th>
                    <th scope="col" class="px-6 py-3">Role</th>
                    <th scope="col" class="px-6 py-3 flex justify-center">
                        Action
                    </th>
                </tr>
            </thead>

            <tbody>
                <tr
                    v-for="i in listUser"
                    :key="i.id"
                    class="bg-white border-b dark:bg-gray-800 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600"
                    @click="setSelectedUser(i.id)"
                >
                    <td class="w-4 p-4">
                        <div class="flex items-center">
                            <input
                                id="checkbox-table-search-1"
                                type="checkbox"
                                class="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 dark:focus:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600"
                            />
                            <label for="checkbox-table-search-1" class="sr-only"
                                >checkbox</label
                            >
                        </div>
                    </td>
                    <th
                        scope="row"
                        class="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white"
                    >
                        {{ i.fullName }}
                    </th>
                    <td class="px-6 py-4">{{ i.createdBy }}</td>
                    <td class="px-6 py-4">{{ i.createdAt }}</td>
                    <td class="px-6 py-4">{{ i.email }}</td>
                    <td class="px-6 py-4">{{ i.phone }}</td>
                    <td class="px-6 py-4">{{ i.roles }}</td>
                    <td class="px-6 py-4 flex justify-center">
                        <button
                            @click="
                                isOpenUpdateUserPopup = !isOpenUpdateUserPopup
                            "
                            type="button"
                            class="h-10 text-green-700 hover:text-white border border-green-700 hover:bg-green-800 focus:ring-4 focus:outline-none focus:ring-green-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center me-2 mb-2 dark:border-green-500 dark:text-green-500 dark:hover:text-white dark:hover:bg-green-600 dark:focus:ring-green-700"
                        >
                            Edit
                        </button>
                        <button
                            @click="isOpenDeletePopup = !isOpenDeletePopup"
                            type="button"
                            class="h-10 text-red-700 hover:text-white border border-red-700 hover:bg-red-800 focus:ring-4 focus:outline-none focus:ring-red-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center me-2 mb-2 dark:border-red-500 dark:text-red-500 dark:hover:text-white dark:hover:bg-red-600 dark:focus:ring-red-900"
                        >
                            Delete
                        </button>
                    </td>
                </tr>
            </tbody>
        </table>
        <!-- <nav
            v-if="false"
            class="flex items-center flex-column flex-wrap md:flex-row justify-between pt-4"
            aria-label="Table navigation"
        >
            <span
                class="text-sm font-normal text-gray-500 dark:text-gray-400 mb-4 md:mb-0 block w-full md:inline md:w-auto"
                >Showing
                <span class="font-semibold text-gray-900 dark:text-white"
                    >1-10</span
                >
                of
                <span class="font-semibold text-gray-900 dark:text-white">{{
                    totalServer
                }}</span></span
            >
            <ul class="inline-flex -space-x-px rtl:space-x-reverse text-sm h-8">
                <li @click="onClickPrevious">
                    <span
                        class="flex items-center justify-center px-3 h-8 ms-0 leading-tight text-gray-500 bg-white border border-gray-300 rounded-s-lg hover:bg-gray-100 hover:text-gray-700 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white"
                    >
                        Previous
                    </span>
                </li>
                <li
                    v-for="p in pigination"
                    :key="p"
                    @click="onClickPage(p)"
                    :class="p === page ? 'bg-blue-300' : ''"
                >
                    <span
                        class="flex items-center justify-center px-3 h-8 leading-tight text-gray-500 border border-gray-300 hover:bg-gray-100 hover:text-gray-700 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white"
                    >
                        {{ p }}
                    </span>
                </li>
                <li @click="onClickNext">
                    <span
                        class="flex items-center justify-center px-3 h-8 leading-tight text-gray-500 bg-white border border-gray-300 rounded-e-lg hover:bg-gray-100 hover:text-gray-700 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white"
                    >
                        Next
                    </span>
                </li>
            </ul>
        </nav> -->
    </div>
    <CreateUserPopup v-model:is-open="isOpenCreateUserPopup"></CreateUserPopup>
    <UpdateUserPopup v-model:is-open="isOpenUpdateUserPopup"></UpdateUserPopup>
    <ComfirmPopup
        :content="confirmDelete.content"
        :title="confirmDelete.title"
        v-model:is-open="isOpenDeletePopup"
        @confirm="onAcceptDelete"
    ></ComfirmPopup>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import CreateUserPopup from "./CreateUserPopup.vue";
import UpdateUserPopup from "./UpdateUserPopup.vue";
import EmptyStates from "@/components/base/EmptyStates.vue";
import ComfirmPopup from "@/components/base/ComfirmPopup.vue";
import { userService } from "@/plugins/axios/user/userService";
import { useUserStore } from "@/stores/userStore";

const isOpenCreateUserPopup = ref(false);
const isOpenUpdateUserPopup = ref(false);
const isOpenDeletePopup = ref(false);

const confirmDelete = ref({
    title: "Warnning delete",
    content: "Select 'Yes' if you are sure you want to delete the user.",
});

const userStore = useUserStore();

const listUser = userStore.listUserComputed;
const selectedUser = userStore.selectedUserComputed;

onMounted(() => {
    userService.getListUser().then((res) => {
        const { data } = res;
        console.log(data);
        userStore.setListUser(data.users);
    });
});

const setSelectedUser = (id: number) => {
    userStore.setSelectedUser(id);
};

const onAcceptDelete = (value: boolean) => {
    if (value && selectedUser.value?.id) {
        userService.deleteUser(selectedUser.value?.id).then(() => {
            userService.getListUser().then((res) => {
                const { data } = res;
                userStore.setListUser(data.users);
            });
        });
    }
};
</script>
