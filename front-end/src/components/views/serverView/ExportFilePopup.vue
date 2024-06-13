<template>
    <TransitionRoot appear :show="isOpen" as="template">
        <Dialog as="div" @close="closeModal" class="relative z-10">
            <TransitionChild as="template" enter="duration-300 ease-out" enter-from="opacity-0" enter-to="opacity-100"
                leave="duration-200 ease-in" leave-from="opacity-100" leave-to="opacity-0">
                <div class="fixed inset-0 bg-black/25" />
            </TransitionChild>

            <div class="fixed inset-0 overflow-y-auto">
                <div class="flex min-h-full items-center justify-center p-4 text-center">
                    <TransitionChild as="template" enter="duration-300 ease-out" enter-from="opacity-0 scale-95"
                        enter-to="opacity-100 scale-100" leave="duration-200 ease-in" leave-from="opacity-100 scale-100"
                        leave-to="opacity-0 scale-95">
                        <DialogPanel
                            class="w-full max-w-md transform overflow-hidden rounded-2xl bg-white p-6 text-left align-middle shadow-xl transition-all">
                            <DialogTitle as="h3" class="text-2xl font-medium leading-6 text-gray-900">
                                Export file
                            </DialogTitle>
                            <form @submit.prevent="onSubmit">

                                <div class="mt-2 mb-2 border-b-2 pb-2">
                                    <h4 class="text-xl font-medium leading-6 text-gray-900">
                                        Pagination
                                    </h4>
                                    <div class="mt-2 flex">
                                        <InputFiled class="w-1/2" label="Page Size" :is-disable="false" :is-required="false"
                                            place-holder="Page size" v-model:model-value="form.pageSize.value
                                                " :errors="form.errors.value.fromPage"></InputFiled>
                                    </div>
                                    <div class="flex">
                                        <InputFiled class="w-1/2 mr-2" label="From page" :is-required="false"
                                            place-holder="Input page" v-model:model-value="form.fromPage.value
                                                " :errors="form.errors.value.fromPage"></InputFiled>
                                        <InputFiled class="w-1/2" label="To page" :is-required="false"
                                            place-holder="Input page" v-model:model-value="form.toPage.value
                                                " :errors="form.errors.value.toPage"></InputFiled>
                                    </div>
                                    <div class="flex">
                                        <SingleSelection class="w-1/2 mr-2" :filter="sort" label="Sort"
                                            v-model:model-value="form.sort.value
                                                " :errors="form.errors.value.sort"></SingleSelection>
                                        <InputFiled class="w-1/2" label="SortBy" :is-required="false"
                                            place-holder="Sort by field name" v-model:model-value="form.sortBy.value
                                                " :errors="form.errors.value.sortBy"></InputFiled>
                                    </div>
                                    <div class="mt-2">
                                        <Switch label="Server status" :is-required="false" v-model:model-value="form.status.value
                                            "></Switch>
                                    </div>
                                </div>
                                <div class="mt-6 flex items-center justify-end gap-x-6">
                                    <button @click="closeModal" type="button"
                                        class="text-sm font-semibold leading-6 text-gray-900">
                                        Cancel
                                    </button>
                                    <button type="submit"
                                        class="rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600">
                                        Export
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
import { TypeSort } from "@/plugins/axios/server/interfaces";
import { ref } from "vue";
import InputFiled from "@/components/base/InputFiled.vue";
import SingleSelection from "@/components/base/SingleSelection.vue";
import type { ISelectionItem } from "@/components/base/interfaces";
import useFormExportServer from "./exportFileForm";
import { useServerStore } from "@/stores/serverStore";

const isOpen = defineModel("isOpen", { type: Boolean });

function closeModal() {
    isOpen.value = false;
}

const sort = ref<Array<ISelectionItem>>([
    {
        key: TypeSort.NONE,
        value: "None",
    },
    {
        key: TypeSort.ASC,
        value: "Asc",
    },
    {
        key: TypeSort.DESC,
        value: "Desc",
    },
]);

const serverStore = useServerStore();
const currentFilter = serverStore.filterServerComputed;

const form = useFormExportServer();

// function onClickClearFilter() {
//     form.createdAtFrom.value = undefined;
//     form.createdAtTo.value = undefined;
//     form.updatedAtFrom.value = undefined;
//     form.updatedAtTo.value = undefined;
//     form.status.value = undefined;
// }

// function onClickApplyCurentFilter() {
//     form.createdAtFrom.value = currentFilter.value?.createdAtFrom
//         ? currentFilter.value?.createdAtFrom.split("T")[0]
//         : undefined;
//     form.createdAtTo.value = currentFilter.value?.createdAtTo
//         ? currentFilter.value?.createdAtTo.split("T")[0]
//         : undefined;
//     form.updatedAtFrom.value = currentFilter.value?.updatedAtFrom
//         ? currentFilter.value?.updatedAtFrom.split("T")[0]
//         : undefined;
//     form.updatedAtTo.value = currentFilter.value?.updatedAtTo
//         ? currentFilter.value?.updatedAtTo.split("T")[0]
//         : undefined;
//     form.status.value = currentFilter.value?.status;
// }

function onSubmit() {
    form.onSubmit().then((res) => {
        closeModal();
    });
}
</script>
