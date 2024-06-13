<template>
    <div class="relative overflow-x-auto sm:rounded-lg">
        <div class="flex flex-column sm:flex-row flex-wrap space-y-4 sm:space-y-0 items-center justify-between pb-4">
            <div class="flex">
                <button @click="onClickCreate" type="button"
                    class="text-white bg-[#1da1f2] hover:bg-[#1da1f2]/90 focus:ring-4 focus:outline-none focus:ring-[#1da1f2]/50 font-medium rounded-lg text-sm px-5 py-2.5 text-center inline-flex items-center dark:focus:ring-[#1da1f2]/55 mx-1">
                    New Server
                </button>
                <button @click="onClickFilter" type="button"
                    class="text-white bg-[#1da1f2] hover:bg-[#1da1f2]/90 focus:ring-4 focus:outline-none focus:ring-[#1da1f2]/50 font-medium rounded-lg text-sm px-5 py-2.5 text-center inline-flex items-center dark:focus:ring-[#1da1f2]/55 mx-1">
                    Filter
                </button>
                <button @click="onClickImport" type="button"
                    class="text-white bg-[#1da1f2] hover:bg-[#1da1f2]/90 focus:ring-4 focus:outline-none focus:ring-[#1da1f2]/50 font-medium rounded-lg text-sm px-5 py-2.5 text-center inline-flex items-center dark:focus:ring-[#1da1f2]/55 mx-1">
                    Import
                </button>
                <button @click="onClickExport" type="button"
                    class="text-white bg-[#1da1f2] hover:bg-[#1da1f2]/90 focus:ring-4 focus:outline-none focus:ring-[#1da1f2]/50 font-medium rounded-lg text-sm px-5 py-2.5 text-center inline-flex items-center dark:focus:ring-[#1da1f2]/55 mx-1">
                    Export
                </button>
            </div>
            <label for="table-search" class="sr-only">Search</label>
            <div class="relative">
                <div
                    class="absolute inset-y-0 left-0 rtl:inset-r-0 rtl:right-0 flex items-center ps-3 pointer-events-none">
                    <svg class="w-5 h-5 text-gray-500 dark:text-gray-400" aria-hidden="true" fill="currentColor"
                        viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
                        <path fill-rule="evenodd"
                            d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z"
                            clip-rule="evenodd"></path>
                    </svg>
                </div>
                <input type="text" id="table-search"
                    class="block p-2 ps-10 text-sm text-gray-900 border border-gray-300 rounded-lg w-80 bg-gray-50 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                    placeholder="Search for items" v-model="query" @input="onChangeSearchQuery" />
            </div>
        </div>
        <EmptyStates v-if="!listServer || !listServer.length"></EmptyStates>
        <table v-else class="w-full text-sm text-left rtl:text-right text-gray-500 dark:text-gray-400">
            <thead class="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400">
                <tr>
                    <th scope="col" class="p-4">
                        <div class="flex items-center">
                            <input id="checkbox-all-search" type="checkbox"
                                class="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 dark:focus:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600" />
                            <label for="checkbox-all-search" class="sr-only">checkbox</label>
                        </div>
                    </th>
                    <th scope="col" class="px-6 py-3">Server Name</th>
                    <th scope="col" class="px-6 py-3">Create By</th>
                    <th scope="col" class="px-6 py-3">Create At</th>
                    <th scope="col" class="px-6 py-3">Last update</th>
                    <th scope="col" class="px-6 py-3">Status</th>
                    <th scope="col" class="px-6 py-3">IPv4</th>
                    <th scope="col" class="px-6 py-3 flex justify-center">
                        Action
                    </th>
                </tr>
            </thead>

            <tbody>
                <tr v-for="server in listServer" :key="server.ID"
                    class="bg-white border-b dark:bg-gray-800 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600">
                    <td class="w-4 p-4">
                        <div class="flex items-center">
                            <input id="checkbox-table-search-1" type="checkbox"
                                class="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 dark:focus:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600" />
                            <label for="checkbox-table-search-1" class="sr-only">checkbox</label>
                        </div>
                    </td>
                    <th scope="row" class="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white">
                        {{ server.name }}
                    </th>
                    <td class="px-6 py-4">{{ server.CreatedAt }}</td>
                    <td class="px-6 py-4">
                        {{ server.CreatedAt }}
                    </td>
                    <td class="px-6 py-4">
                        {{ server.UpdatedAt }}
                    </td>
                    <td class="px-6 py-4">
                        {{ server.status === true ? "ON" : "OFF" }}
                    </td>
                    <td class="px-6 py-4">{{ server.ip }}</td>
                    <td class="px-6 py-4 flex justify-center">
                        <button @click="onClickEdit(server.ID)" type="button"
                            class="h-10 text-green-700 hover:text-white border border-green-700 hover:bg-green-800 focus:ring-4 focus:outline-none focus:ring-green-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center me-2 mb-2 dark:border-green-500 dark:text-green-500 dark:hover:text-white dark:hover:bg-green-600 dark:focus:ring-green-700">
                            Edit
                        </button>
                        <button @click="onClickDelete(server.ID)" type="button"
                            class="h-10 text-red-700 hover:text-white border border-red-700 hover:bg-red-800 focus:ring-4 focus:outline-none focus:ring-red-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center me-2 mb-2 dark:border-red-500 dark:text-red-500 dark:hover:text-white dark:hover:bg-red-600 dark:focus:ring-red-900">
                            Delete
                        </button>
                    </td>
                </tr>
            </tbody>
        </table>
        <nav class="flex items-center flex-column flex-wrap md:flex-row justify-between pt-4"
            aria-label="Table navigation">
            <span
                class="text-sm font-normal text-gray-500 dark:text-gray-400 mb-4 md:mb-0 block w-full md:inline md:w-auto">Showing
                <span class="font-semibold text-gray-900 dark:text-white">1-10</span>
                of
                <span class="font-semibold text-gray-900 dark:text-white">{{
                    totalServer
                    }}</span></span>
            <ul class="inline-flex -space-x-px rtl:space-x-reverse text-sm h-8">
                <li @click="onClickPrevious">
                    <span
                        class="flex items-center justify-center px-3 h-8 ms-0 leading-tight text-gray-500 bg-white border border-gray-300 rounded-s-lg hover:bg-gray-100 hover:text-gray-700 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white">
                        Previous
                    </span>
                </li>
                <li v-for="p in pigination" :key="p" @click="onClickPage(p)" :class="p === page ? 'bg-blue-300' : ''">
                    <span
                        class="flex items-center justify-center px-3 h-8 leading-tight text-gray-500 border border-gray-300 hover:bg-gray-100 hover:text-gray-700 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white">
                        {{ p }}
                    </span>
                </li>
                <li @click="onClickNext">
                    <span
                        class="flex items-center justify-center px-3 h-8 leading-tight text-gray-500 bg-white border border-gray-300 rounded-e-lg hover:bg-gray-100 hover:text-gray-700 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white">
                        Next
                    </span>
                </li>
            </ul>
        </nav>
    </div>
    <EditServerPopup v-model:isOpen="isOpenEditPopup"></EditServerPopup>
    <CreateServerPopup v-model:isOpen="isOpenCreatePopup"></CreateServerPopup>
    <FilterServerPopup v-model:is-open="isOpenFilterPopup"></FilterServerPopup>
    <ImportFilePopup v-model:is-open="isOpenImportPopup"></ImportFilePopup>
    <ExportFilePopup v-model:is-open="isOpenExportPopup"></ExportFilePopup>
    <ComfirmPopup :content="confirmDelete.content" :title="confirmDelete.title" v-model:is-open="isOpenDeletePopup"
        @confirm="onAcceptDelete"></ComfirmPopup>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { useServerStore } from "@/stores/serverStore";
import { Status } from "./../interfaces";
import EditServerPopup from "./EditServerPopup.vue";
import CreateServerPopup from "./CreateServerPopup.vue";
import FilterServerPopup from "./FilterServerPopup.vue";
import ImportFilePopup from "./ImportFilePopup.vue";
import ExportFilePopup from "./ExportFilePopup.vue";
import EmptyStates from "@/components/base/EmptyStates.vue";
import { DefaultQuery } from "../constants";
import ComfirmPopup from "@/components/base/ComfirmPopup.vue";
import { serverService } from "@/plugins/axios/server/serverService";
import type {
    TypeSort,
    IListServerRequest
} from "@/plugins/axios/server/interfaces";
import { useMainStore } from "@/stores/mainStore";

const serverStore = useServerStore();
const mainStore = useMainStore();

const selectedServer = serverStore.selectedServerComputed;

const page = ref(0);

const confirmDelete = ref({
    title: "Warnning delete",
    content: "Select 'Yes' if you are sure you want to delete the server.",
});

const pigination = computed(() => {
    const result: number[] = [];
    for (let index: number = page.value - 2; index < page.value + 2; index++) {
        if (index * 10 >= 0 && (index + 1) * 10 - 9 <= totalServer.value) {
            result.push(index);
        }
    }
    return result;
});

const getListServer = (req: IListServerRequest) => {
    serverService.getListServer(req).then((response) => {
        const { data } = response;
        console.log(data);
        
        serverStore.updateServers(data.data);
        serverStore.updateTotalServer(data.total);
    });
};

onMounted(() => {
    getListServer(DefaultQuery);
});

const isOpenEditPopup = ref(false);
const isOpenCreatePopup = ref(false);
const isOpenDeletePopup = ref(false);
const isOpenFilterPopup = ref(false);
const isOpenImportPopup = ref(false);
const isOpenExportPopup = ref(false);

const listServer = serverStore.servers;
const totalServer = serverStore.total;
const filterServer = serverStore.filterServerComputed;

const onClickDelete = (id: number) => {
    serverStore.setSelectedServer(id);
    isOpenDeletePopup.value = !isOpenDeletePopup.value;
};

const onClickEdit = (id: number) => {
    serverStore.setSelectedServer(id);
    isOpenEditPopup.value = !isOpenEditPopup.value;
};

const onClickNext = () => {
    if (totalServer && page.value + 1 < totalServer.value / 10) {
        page.value++;
        getListServer({
            limit: 10,
            offset: page.value * 10,
            field: "created_at",
        });
    }
};

const onClickPrevious = () => {
    if (page.value > 0) {
        page.value--;
        getListServer({
            limit: 10,
            offset: page.value * 10,
            field: "created_at",
        });
    }
};

function onClickPage(p: number) {
    page.value = p;
    getListServer({
        limit: 10,
        offset: page.value * 10,
        field: "created_at",
    });
}

const onClickCreate = () => {
    isOpenCreatePopup.value = !isOpenCreatePopup.value;
};

const onClickFilter = () => {
    isOpenFilterPopup.value = !isOpenFilterPopup.value;
};

const onClickImport = () => {
    isOpenImportPopup.value = !isOpenImportPopup.value;
};

const onClickExport = () => {
    isOpenExportPopup.value = !isOpenExportPopup.value;
};

const onAcceptDelete = (value: boolean) => {
    if (value && selectedServer.value?.ID) {
        serverService
            .deleteServer(selectedServer.value?.ID)
            .then((response) => {
                page.value = 0;
                mainStore.showNofitication({
                    typeNotification: "infor",
                    title: "Delete server",
                    content: `Delete server successfuly`,
                });
                getListServer(DefaultQuery);
            });
    }
};

const query = ref("");
let delaySearch: any;
function onChangeSearchQuery(e: Event) {
    clearTimeout(delaySearch);
    delaySearch = setTimeout(() => {
        getListServer(DefaultQuery);
    }, 1000);
}
</script>
