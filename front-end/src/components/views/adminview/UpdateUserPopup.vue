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
                                Update User
                            </DialogTitle>

                            <div class="mt-2">
                                <form @submit.prevent="">
                                    <div class="space-y-12">
                                        <div
                                            class="border-b border-gray-900/10 pb-12"
                                        >
                                            <h2
                                                class="text-base font-semibold leading-7 text-gray-900"
                                            >
                                                User profile
                                            </h2>
                                            <p
                                                class="mt-1 text-sm leading-6 text-gray-600"
                                            >
                                                This information will be
                                                displayed when save
                                            </p>

                                            <div class="mt-10">
                                                <InputFiled
                                                    :is-required="true"
                                                    label="User email"
                                                    v-model:model-value="
                                                        form.email.value
                                                    "
                                                    :errors="
                                                        form.errors.value.email
                                                    "
                                                ></InputFiled>
                                                <InputFiled
                                                    :is-required="true"
                                                    label="Fullname"
                                                    v-model:model-value="
                                                        form.fullName.value
                                                    "
                                                    :errors="
                                                        form.errors.value
                                                            .fullName
                                                    "
                                                ></InputFiled>
                                                <InputFiled
                                                    :is-required="true"
                                                    label="Phone nunber"
                                                    v-model:model-value="
                                                        form.phone.value
                                                    "
                                                    :errors="
                                                        form.errors.value.phone
                                                    "
                                                ></InputFiled>
                                                <SingleSelection
                                                    label="Role"
                                                    :is-required="true"
                                                    :filter="[]"
                                                    v-model:model-value="
                                                        form.roles.value
                                                    "
                                                    :errors="
                                                        form.errors.value.roles
                                                    "
                                                ></SingleSelection>
                                            </div>
                                        </div>
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
                                            type="submit"
                                            class="rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
                                        >
                                            Create
                                        </button>
                                    </div>
                                </form>
                            </div>
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
import InputFiled from "@/components/base/InputFiled.vue";
import SingleSelection from "@/components/base/SingleSelection.vue";
import useUpdateUserForm from "./updateUserForm";
import { useUserStore } from "@/stores/userStore";
import { onMounted } from "vue";
import { userService } from "@/plugins/axios/user/userService";

const isOpen = defineModel("isOpen", { type: Boolean });

function closeModal() {
    isOpen.value = false;
}
const form = useUpdateUserForm();

const userStore = useUserStore();
const selectedUser = userStore.selectedUserComputed;

onMounted(() => {
    form.setFieldValue("email", selectedUser.value?.email);
    form.setFieldValue("fullName", selectedUser.value?.fullName);
    form.setFieldValue("phone", selectedUser.value?.phone);
    form.setFieldValue(
        "roles",
        selectedUser.value?.roles === "RoleAdmin" ? "Admin" : "User"
    );
});

const onSubmit = () => {
    form.onSubmit().then(() => {
        userService.getListUser().then((res) => {
            const { data } = res;
            userStore.setListUser(data.users);
            closeModal();
        });
    });
};
</script>
