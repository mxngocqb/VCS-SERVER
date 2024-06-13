<template>
    <form @submit.prevent="onSubmit">
        <div class="space-y-12">
            <div class="border-b border-gray-900/10 pb-12">
                <h2 class="text-base font-semibold leading-7 text-gray-900">
                    Server profile
                </h2>
                <p class="mt-1 text-sm leading-6 text-gray-600">
                    This information will be displayed when save
                </p>

                <div class="mt-10">
                    <InputFiled
                        :is-required="true"
                        label="Server name"
                        v-model:model-value="form.name.value"
                    ></InputFiled>
                    <InputFiled
                        :is-required="true"
                        label="Server IPv4"
                        v-model:model-value="form.ipv4.value"
                        :errors="form.errors.value.ipv4"
                    ></InputFiled>
                </div>
            </div>
        </div>

        <div class="mt-6 flex items-center justify-end gap-x-6">
            <button
                type="button"
                class="text-sm font-semibold leading-6 text-gray-900"
                @click="emit('closeForm')"
            >
                Cancel
            </button>
            <button
                type="submit"
                class="rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
            >
                Save
            </button>
        </div>
    </form>
</template>

<script setup lang="ts">
import useUpdateServerForm from "./updateServerFrom";
import InputFiled from "@/components/base/InputFiled.vue";
import { useServerStore } from "@/stores/serverStore";
import { onMounted, onUnmounted } from "vue";
import { Status } from "../interfaces";
import { serverService } from "@/plugins/axios/server/serverService";
import type { IListServerRequest } from "@/plugins/axios/server/interfaces";
import { DefaultQuery } from "../constants";

const serverStore = useServerStore();

const selectedServer = serverStore.selectedServerComputed;
const filterServer = serverStore.filterServerComputed;

const form = useUpdateServerForm();

onMounted(() => {
    form.resetForm();
    form.setFieldValue("name", selectedServer.value?.name);
    form.setFieldValue("ipv4", selectedServer.value?.ip);
});

onUnmounted(() => {
    form.resetForm();
});

const emit = defineEmits(["closeForm"]);

const getListServer = (req: IListServerRequest) => {
    serverService.getListServer(req).then((response) => {
        const { data } = response;
        serverStore.updateServers(data.data);
        serverStore.updateTotalServer(data.total);
    });
};

const onSubmit = async () => {
    const result = await form.onSubmit();
    if (result) {
        getListServer(DefaultQuery);
        emit("closeForm");
    }
};
</script>
