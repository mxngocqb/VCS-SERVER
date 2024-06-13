<template>
    <TransitionRoot appear :show="isOpen" as="template">
        <Dialog as="div" @close="closeModal" class="relative z-10">
            <TransitionChild
                as="template"
                enter="duration-300 ease-out"
                enter-from="opacity-0"
                enter-to="opacity-100"
                leave="duration-200 ease-in"
                leave-from="opacity-100"
                leave-to="opacity-0"
            >
                <div class="fixed inset-0 bg-black/25" />
            </TransitionChild>

            <div class="fixed inset-0 overflow-y-auto">
                <div
                    class="flex min-h-full items-center justify-center p-4 text-center"
                >
                    <TransitionChild
                        as="template"
                        enter="duration-300 ease-out"
                        enter-from="opacity-0 scale-95"
                        enter-to="opacity-100 scale-100"
                        leave="duration-200 ease-in"
                        leave-from="opacity-100 scale-100"
                        leave-to="opacity-0 scale-95"
                    >
                        <DialogPanel
                            class="w-full max-w-md transform overflow-hidden rounded-2xl bg-white p-6 text-left align-middle shadow-xl transition-all"
                        >
                            <DialogTitle
                                as="h3"
                                class="text-2xl font-medium leading-6 text-gray-900"
                            >
                                Filter
                            </DialogTitle>
                            <form @submit.prevent="onSubmit">
                                <div class="mt-2">
                                    <h5
                                        class="text-l font-medium leading-6 text-gray-900"
                                    >
                                        Created At
                                    </h5>
                                    <div class="flex justify-between">
                                        <DatePicker
                                            class="px-4"
                                            label="From"
                                            :is-required="false"
                                            v-model:value="
                                                form.createdAtFrom.value
                                            "
                                            :errors="
                                                form.errors.value.createdAtFrom
                                            "
                                        ></DatePicker>
                                        <DatePicker
                                            class="px-4"
                                            label="To"
                                            :is-required="false"
                                            v-model:value="
                                                form.createdAtTo.value
                                            "
                                            :errors="
                                                form.errors.value.createdAtTo
                                            "
                                        ></DatePicker>
                                    </div>
                                </div>
                                <div class="mt-2">
                                    <h5
                                        class="text-l font-medium leading-6 text-gray-900"
                                    >
                                        Last update
                                    </h5>
                                    <div class="flex justify-between">
                                        <DatePicker
                                            class="px-4"
                                            label="From"
                                            :is-required="false"
                                            :errors="
                                                form.errors.value.updatedAtFrom
                                            "
                                            v-model:value="
                                                form.updatedAtFrom.value
                                            "
                                        ></DatePicker>
                                        <DatePicker
                                            class="px-4"
                                            label="To"
                                            :is-required="false"
                                            :errors="
                                                form.errors.value.updatedAtTo
                                            "
                                            v-model:value="
                                                form.updatedAtTo.value
                                            "
                                        ></DatePicker>
                                    </div>
                                </div>
                                <div class="mt-2">
                                    <Switch
                                        label="Server status"
                                        :is-required="false"
                                        v-model:model-value="form.status.value"
                                    ></Switch>
                                </div>
                                <div
                                    class="mt-6 flex items-center justify-end gap-x-6"
                                >
                                    <button
                                        type="button"
                                        class="text-sm font-semibold leading-6 text-gray-900"
                                        @click="closeModal"
                                    >
                                        Cancel
                                    </button>
                                    <button
                                        @click="onClear"
                                        type="button"
                                        class="rounded-md bg-orange-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-orange-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-orange-600"
                                    >
                                        Clear
                                    </button>
                                    <button
                                        type="submit"
                                        class="rounded-md bg-blue-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-blue-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600"
                                    >
                                        Apply
                                    </button>
                                </div>
                            </form>
                        </DialogPanel>
                    </TransitionChild>
                </div>
            </div>
        </Dialog>
    </TransitionRoot>
</template>

<script setup lang="ts">
import {
    TransitionRoot,
    TransitionChild,
    Dialog,
    DialogPanel,
    DialogTitle,
} from "@headlessui/vue";
import DatePicker from "@/components/base/DatePicker.vue";
import Switch from "@/components/base/Switch.vue";
import useFilterForm from "./filterServerFrom";
import { useServerStore } from "@/stores/serverStore";
import {
    ServerStatus,
    type IListServerRequest,
} from "@/plugins/axios/server/interfaces";
import { serverService } from "@/plugins/axios/server/serverService";
import { DefaultQuery } from "../constants";
import { onMounted } from "vue";

const isOpen = defineModel("isOpen", { type: Boolean });

function closeModal() {
    isOpen.value = false;
}

const form = useFilterForm();

const serverStore = useServerStore();
const filterServer = serverStore.filterServerComputed;

function onClear() {
    form.resetForm();
}

const getListServer = (req: IListServerRequest) => {
    serverService.getListServer(req).then((response) => {
        const { data } = response;
        serverStore.updateServers(data.data);
        serverStore.updateTotalServer(data.total);
    });
};

onMounted(() => {
    form.resetForm();
    form.setFieldValue("createdAtFrom", filterServer.value?.createdAtFrom);
    form.setFieldValue("createdAtTo", filterServer.value?.createdAtTo);
    form.setFieldValue("updatedAtFrom", filterServer.value?.updatedAtFrom);
    form.setFieldValue("updatedAtTo", filterServer.value?.updatedAtTo);
    form.setFieldValue(
        "status",
        filterServer.value?.status
            ? filterServer.value?.status === ServerStatus.ON
            : undefined
    );
});

async function onSubmit() {
    const result = await form.onSubmit();
    if (result) {
        await getListServer(DefaultQuery);
        closeModal();
    }
}
</script>
