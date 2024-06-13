<template>
    <RouterView :key="$route.path" />
    <ul class="fixed top-2 right-0">
        <li v-for="nofi in nofitications">
            <NofiticationItem
                :content="nofi?.content"
                :title="nofi?.title"
                :type-notification="nofi?.typeNotification"
                class="my-2 mx-2"
            ></NofiticationItem>
        </li>
    </ul>
</template>

<script setup lang="ts">
import { useMainStore } from "./stores/mainStore";
import NofiticationItem from "@/components/base/NofiticationItem.vue";
const mainStore = useMainStore();

const nofitications = mainStore.nofiticationsComputed;
</script>

<style scoped></style>
