<template>
    <div class="flex flex-column sm:flex-row flex-wrap space-y-4 sm:space-y-0 items-center justify-between pb-4">
        <div class="flex">
            <button @click="isOpenSendReportPopup = !isOpenSendReportPopup" type="button"
                class="text-white bg-[#1da1f2] hover:bg-[#1da1f2]/90 focus:ring-4 focus:outline-none focus:ring-[#1da1f2]/50 font-medium rounded-lg text-sm px-5 py-2.5 text-center inline-flex items-center dark:focus:ring-[#1da1f2]/55 mx-1">
                Send report
            </button>
        </div>
    </div>
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-6">
        <div class="bg-white rounded-md border border-gray-100 p-6 shadow-md shadow-black/5">
            <div class="flex justify-between mb-6">
                <div>
                    <div class="flex items-center mb-1">
                        <div class="text-2xl font-semibold">{{ totalServers }}</div>
                    </div>
                    <div class="text-sm font-medium text-gray-400">
                        Total server
                    </div>
                </div>
            </div>

            <RouterLink :to="{ name: 'server' }" class="text-[#f84525] font-medium text-sm hover:text-red-800">View
            </RouterLink>
        </div>
        <div class="bg-white rounded-md border border-gray-100 p-6 shadow-md shadow-black/5">
            <div class="flex justify-between mb-6">
                <div>
                    <div class="flex items-center mb-1">
                        <div class="text-2xl font-semibold">{{ serverOn }}</div>
                    </div>
                    <div class="text-sm font-medium text-gray-400">
                        Server On
                    </div>
                </div>
            </div>

            <RouterLink :to="{ name: 'server' }" class="text-[#f84525] font-medium text-sm hover:text-red-800">View
            </RouterLink>
        </div>
        <div class="bg-white rounded-md border border-gray-100 p-6 shadow-md shadow-black/5">
            <div class="flex justify-between mb-6">
                <div>
                    <div class="flex items-center mb-1">
                        <div class="text-2xl font-semibold">{{ serverOff }}</div>
                    </div>
                    <div class="text-sm font-medium text-gray-400">
                        Server off
                    </div>
                </div>
            </div>

            <RouterLink :to="{ name: 'server' }" class="text-[#f84525] font-medium text-sm hover:text-red-800">View
            </RouterLink>
        </div>
    </div>
    <SendReportPopup v-model:is-open="isOpenSendReportPopup"></SendReportPopup>
</template>
<script setup lang="ts">
import { ref, onMounted } from 'vue';
import SendReportPopup from './SendReportPopup.vue';
import axios from 'axios';
import { serverService } from "@/plugins/axios/server/serverService";

const isOpenSendReportPopup = ref(false);

const totalServers = ref(0);
const serverOn = ref(0);
const serverOff = ref(0);


const updateServerStatus = () => {
    serverService.getListServerStatus().
        then(response => {
            const data = response.data as IServerStatusResponse;
            console.log(data);

            totalServers.value = data.online + data.offline;
            serverOn.value = data.online;
            serverOff.value = data.offline;
        }).catch(error => {
            console.error('Failed to fetch server status:', error);
        });
};

onMounted(() => {
    updateServerStatus();
});
</script>

<script lang="ts">
export interface IServerStatusResponse {
    online: number;
    offline: number;
}
</script>
