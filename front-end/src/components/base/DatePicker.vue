<template>
    <div class="relative max-w-sm">
        <label
            for="input"
            class="block mb-2 text-sm font-bold text-gray-900 dark:text-gray-400"
            >{{ props.label }}
            <span v-if="isRequired" class="text-red-800">*</span></label
        >

        <input
            @focusin="onFocusIn"
            @focusout="onFocusOut"
            :type="isFocus ? 'date' : 'text'"
            class="relative bg-gray-50 border focus:text-gray-900 border-gray-300 text-gray-400 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full ps-10 p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
            v-model="modelValue"
            :placeholder="props.placeHolder"
        />
        <p
            v-if="props.errors"
            class="mt-2 text-xs italic text-red-600 dark:text-red-500"
        >
            <span class="font-medium">! {{ props.errors }}</span>
        </p>
    </div>
</template>

<script setup lang="ts">
import { ref } from "vue";

const props = defineProps<{
    label: string;
    isRequired: boolean;
    placeHolder?: string;
    errors?: string;
}>();
const modelValue = defineModel<string>("value");
const isFocus = ref(false);
function onFocusIn(e: FocusEvent) {
    isFocus.value = true;
}

function onFocusOut(e: FocusEvent) {
    isFocus.value = false;
}
</script>

<style scoped>
[type="date"]::-webkit-calendar-picker-indicator {
    opacity: 1;
    left: 0.7rem;
    position: absolute;
    top: 0.7rem;
    width: 100%;
}
</style>
