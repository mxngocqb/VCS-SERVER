<template>
    <div class="relative overflow-x-auto sm:rounded-lg sm:boder-2">
        <div
            class="w-full uppercase font-bold text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400 px-6 py-3"
        >
            Uploaded files
        </div>
        <EmptyStates v-if="!fileList.length"></EmptyStates>
        <table
            v-else
            class="w-full text-sm text-left rtl:text-right text-gray-500 dark:text-gray-400"
        >
            <tbody>
                <tr
                    v-for="f in fileList"
                    :key="f.Id"
                    class="bg-white border-b dark:bg-gray-800 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600"
                >
                    <td class="w-4 h-4">
                        <IconExcell />
                    </td>
                    <th scope="col" class="px-6 py-3 w-1/2">
                        {{ f.fileName }}
                    </th>
                    <th scope="col" class="px-6 py-3">{{ f.fileSize }}</th>
                    <th scope="col" class="px-6 py-3">
                        <div>{{ f.status }}</div>
                    </th>
                    <td>
                        <span
                            @click="onClickDownload(f.filePath)"
                            class="px-6 py-3 flex justify-center text-blue-500 hover:underline hover:cursor-pointer"
                            >Download</span
                        >
                    </td>
                </tr>
            </tbody>
        </table>
    </div>
</template>
<script setup lang="ts">
import EmptyStates from "@/components/base/EmptyStates.vue";
import IconExcell from "@/components/icons/IconExcell.vue";
import { fileService } from "@/plugins/axios/storage/fileService";
import type { IFile } from "@/plugins/axios/storage/interfaces";
import { useUserStore } from "@/stores/userStore";
import { onMounted, ref } from "vue";

const fileList = ref<Array<IFile>>([]);
const userStore = useUserStore();

onMounted(() => {
    if (currentUser.value?.id) {
        fileService.getAllFileOfUser(currentUser.value?.id).then((res) => {
            const { data } = res;
            fileList.value = data.files;
        });
    }
});

const currentUser = userStore.currentUser;

const onClickDownload = (path: string) => {
    fileService.download(path).then((response) => {
        const href = URL.createObjectURL(response.data);

        // create "a" HTML element with href to file & click
        const link = document.createElement("a");
        link.href = href;
        link.setAttribute("download", "file.xlsx"); //or any other extension
        document.body.appendChild(link);
        link.click();

        // clean up "a" element & remove ObjectURL
        document.body.removeChild(link);
        URL.revokeObjectURL(href);
    });
};
</script>
