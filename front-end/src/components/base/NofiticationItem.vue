<template>
    <div :class="styleAlert" class="w-96 rounded-lg" role="alert">
        <p class="font-bold">{{ props.title }}</p>
        <p>{{ props.content }}</p>
    </div>
</template>
<script setup lang="ts">
import { onMounted, ref } from "vue";

const props = defineProps<{
    title: string;
    content: string;
    typeNotification: "error" | "infor" | "warnning";
}>();

const warnningStyle =
    "bg-orange-100 border-l-4 border-b border-orange-500 text-orange-700 p-4";
const errorStyle =
    "bg-red-100 border-l-4 border-b border-red-500 text-red-700 p-4";
const inforStyle =
    "bg-blue-100 border-t border-l-4 border-b border-blue-500 text-blue-700 px-4 ";

const styleAlert = ref("infor");
onMounted(() => {
    if (props.typeNotification === "error") {
        styleAlert.value = errorStyle;
    } else if (props.typeNotification === "infor") {
        styleAlert.value = inforStyle;
    } else styleAlert.value = warnningStyle;
    setTimeout(() => {
        styleAlert.value +=
            "transition duration-500 ease-in-out opacity-0 translate-x-full";
    }, 2500);
});
</script>
