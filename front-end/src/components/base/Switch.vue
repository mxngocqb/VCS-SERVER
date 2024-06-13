<template>
    <div class="mb-4">
        <label
            for="input"
            class="block mb-2 text-sm font-bold text-gray-900 dark:text-gray-400"
            >{{ props.label }}
            <span v-if="isRequired" class="text-red-800">*</span></label
        >
        <Switch
            v-model="value"
            :class="value ? 'bg-blue-600' : 'bg-gray-200'"
            class="relative inline-flex h-[38px] w-[74px] shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus-visible:ring-2 focus-visible:ring-white/75"
        >
            <span class="sr-only">Use setting</span>
            <span
                aria-hidden="true"
                :class="value ? 'translate-x-9' : 'translate-x-0'"
                class="pointer-events-none inline-block h-[34px] w-[34px] transform rounded-full bg-white shadow-lg ring-0 transition duration-200 ease-in-out"
            />
        </Switch>
        <p
            v-if="props.errors"
            class="mt-2 text-sm text-red-600 dark:text-red-500"
        >
            <span class="font-medium">{{ props.errors }}</span>
        </p>
    </div>
</template>
<script setup lang="ts">
import { Switch } from "@headlessui/vue";

const props = defineProps<{
    label: string;
    isRequired: boolean;
    errors?: string;
}>();

const value = defineModel({ type: Boolean });
</script>
