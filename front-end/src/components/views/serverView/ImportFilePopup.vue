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
                                class="text-2xl font-medium leading-6 text-gray-900 my-2"
                            >
                                Import a file
                            </DialogTitle>

                            <div
                                class="flex items-center justify-center w-full"
                            >
                                <label
                                    for="dropzone-file"
                                    class="flex flex-col items-center justify-center w-full h-64 border-2 border-gray-300 border-dashed rounded-lg cursor-pointer bg-gray-50 dark:hover:bg-bray-800 dark:bg-gray-700 hover:bg-gray-100 dark:border-gray-600 dark:hover:border-gray-500 dark:hover:bg-gray-600"
                                >
                                    <div
                                        class="flex flex-col items-center justify-center pt-5 pb-6"
                                    >
                                        <svg
                                            class="w-8 h-8 mb-4 text-gray-500 dark:text-gray-400"
                                            aria-hidden="true"
                                            xmlns="http://www.w3.org/2000/svg"
                                            fill="none"
                                            viewBox="0 0 20 16"
                                        >
                                            <path
                                                stroke="currentColor"
                                                stroke-linecap="round"
                                                stroke-linejoin="round"
                                                stroke-width="2"
                                                d="M13 13h3a3 3 0 0 0 0-6h-.025A5.56 5.56 0 0 0 16 6.5 5.5 5.5 0 0 0 5.207 5.021C5.137 5.017 5.071 5 5 5a4 4 0 0 0 0 8h2.167M10 15V6m0 0L8 8m2-2 2 2"
                                            />
                                        </svg>
                                        <p
                                            class="mb-2 text-sm text-gray-500 dark:text-gray-400"
                                        >
                                            <span class="font-semibold"
                                                >Click to upload</span
                                            >
                                            or drag and drop
                                        </p>
                                        <p
                                            class="text-xs text-gray-500 dark:text-gray-400"
                                        >
                                            SVG, PNG, JPG or GIF (MAX.
                                            800x400px)
                                        </p>
                                    </div>
                                    <input
                                        id="dropzone-file"
                                        type="file"
                                        class="hidden"
                                        @change="onChange"
                                    />
                                </label>
                            </div>
                            <p class="text-sm py-2">File: {{ file?.name }}</p>
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
                                    @click="onSubmit"
                                    class="rounded-md bg-blue-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-blue-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600"
                                >
                                    Import
                                </button>
                            </div>
                        </DialogPanel>
                    </TransitionChild>
                </div>
            </div>
        </Dialog>
    </TransitionRoot>
</template>
<script setup lang="ts">
import { serverService } from "@/plugins/axios/server/serverService";
import { useMainStore } from "@/stores/mainStore";
import {
    TransitionRoot,
    TransitionChild,
    Dialog,
    DialogPanel,
    DialogTitle,
} from "@headlessui/vue";
import { ref } from "vue";

const isOpen = defineModel("isOpen", { type: Boolean });

function closeModal() {
    isOpen.value = false;
}

const file = ref<File>();

const mainStore = useMainStore();

function onSubmit() {
    if (file.value) {
        if (
            file.value.type !==
            "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
        ) {
            mainStore.showNofitication({
                typeNotification: "error",
                title: "Import file",
                content: "Invalid type",
            });
            return;
        }
        if (file.value.size > 10000000) {
            mainStore.showNofitication({
                typeNotification: "error",
                title: "Import file",
                content: "File size lagger 10Mb",
            });
            return;
        }
        
        serverService.importServer(file.value).then((res) => {
            closeModal();
        });
    }
}

function onChange(e: Event) {
    if (e.target instanceof HTMLInputElement && e.target?.files) {
        file.value = e.target.files[0] || undefined;
    }
}
</script>
