<template>
    <div class="mb-4">
        <label
            for="input"
            class="block mb-2 text-sm font-bold text-gray-900 dark:text-gray-400"
            >{{ props.label }}
            <span v-if="isRequired" class="text-red-800">*</span></label
        >
        <input
            :type="props.type"
            v-model="modelValue"
            :required="isRequired"
            class="bg-gray-50 border focus:text-gray-900 border-gray-300 text-gray-400 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
            :placeholder="props.placeHolder"
            :disabled="props.isDisable"
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
import type { InputTypeHTMLAttribute } from "vue";

const props = defineProps<{
    label: string;
    isRequired: boolean;
    placeHolder?: string;
    errors?: string;
    isDisable?: boolean;
    type?: InputTypeHTMLAttribute | "text";
}>();

const modelValue = defineModel({ type: String });
</script>
