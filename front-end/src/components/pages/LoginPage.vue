<template>
    <div class="w-full h-full bg-gray-50 dark:bg-gray-900">
        <div
            class="flex flex-col items-center justify-center px-6 py-8 mx-auto h-full lg:py-0"
        >
            <a
                href="#"
                class="flex items-center mb-6 text-2xl font-semibold text-gray-900 dark:text-white"
            >
                <img class="w-8 h-8 mr-2" :src="logo" alt="logo" />
                Server Management
            </a>
            <div
                class="w-full bg-white rounded-lg shadow dark:border md:mt-0 sm:max-w-md xl:p-0 dark:bg-gray-800 dark:border-gray-700"
            >
                <div class="p-6 space-y-4 md:space-y-6 sm:p-8">
                    <h1
                        class="text-xl font-bold leading-tight tracking-tight text-gray-900 md:text-2xl dark:text-white"
                    >
                        Sign in to your account
                    </h1>
                    <form
                        class="space-y-4 md:space-y-6"
                        @submit.prevent="onSubmit"
                    >
                        <div>
                            <label
                                for="username"
                                class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
                                >Your username</label
                            >
                            <input
                                type="text"
                                v-model="username"
                                v-bind="usernameAttrs"
                                id="username"
                                class="bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-primary-600 focus:border-primary-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                                placeholder="yourusername"
                                required="true"
                            />
                            <p
                                v-if="errors.username"
                                class="mt-2 text-xs italic text-red-600 dark:text-red-500"
                            >
                                <span class="font-medium"
                                    >! {{ errors.username }}</span
                                >
                            </p>
                        </div>
                        <div>
                            <label
                                for="password"
                                class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
                                >Password</label
                            >
                            <input
                                type="password"
                                v-model="password"
                                v-bind="passwordAttrs"
                                id="password"
                                placeholder="••••••••"
                                class="bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-primary-600 focus:border-primary-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                                required="true"
                            />
                            <p
                                v-if="errors.password"
                                class="mt-2 text-xs italic text-red-600 dark:text-red-500"
                            >
                                <span class="font-medium"
                                    >! {{ errors.password }}</span
                                >
                            </p>
                        </div>
                        <div class="flex items-center justify-between">
                            <div class="flex items-start">
                                <div class="flex items-center h-5">
                                    <input
                                        id="remember"
                                        aria-describedby="remember"
                                        type="checkbox"
                                        class="w-4 h-4 border border-gray-300 rounded bg-gray-50 focus:ring-3 focus:ring-primary-300 dark:bg-gray-700 dark:border-gray-600 dark:focus:ring-primary-600 dark:ring-offset-gray-800"
                                    />
                                </div>
                                <div class="ml-3 text-sm">
                                    <label
                                        for="remember"
                                        class="text-gray-500 dark:text-gray-300"
                                        >Remember me</label
                                    >
                                </div>
                            </div>
                            <a
                                href="#"
                                class="text-sm font-medium text-primary-600 hover:underline dark:text-primary-500"
                                >Forgot password?</a
                            >
                        </div>
                        <div>
                            <button
                                type="submit"
                                class="flex w-full justify-center rounded-md bg-blue-500 px-3 py-1.5 text-sm font-semibold leading-6 text-white shadow-sm hover:bg-blue-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
                            >
                                Sign in
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import logo from "@/assets/_995ed77b-bb47-45ba-aa62-42c8bc67e68a-removebg-preview.png";
import { authService } from "@/plugins/axios/auth/authService";
import { userService } from "@/plugins/axios/user/userService";
import { SessionStorageKey } from "@/stores/constants";
import { useUserStore } from "@/stores/userStore";
import { useForm } from "vee-validate";
import { useRouter } from "vue-router";
import * as yup from "yup";

const { values, errors, defineField, handleSubmit } = useForm({
    validationSchema: yup.object({
        username: yup.string().required(),
        password: yup.string().required(),
    }),
});

const [username, usernameAttrs] = defineField("username");
const [password, passwordAttrs] = defineField("password");
const routes = useRouter();
const userStore = useUserStore();

const onSubmit = handleSubmit(async (values) => {
    authService
        .login({
            username: values["username"],
            password: values["password"],
        })
        .then((response) => {
            const { data } = response;
            sessionStorage.setItem(
                SessionStorageKey.AUTH_TOKEN,
                data?.token
            );
            sessionStorage.setItem(
                SessionStorageKey.AUTH_REFRESH_TOKEN,
                data?.refreshToken
            );
            userService
                .getUserByUsername({ username: values["username"] }) // Assuming you have a getUserByUsername function
                .then((response) => {
                    const { data } = response;
                    console.log(data);
                    
                    userStore.setUser(data);
                });

            routes.push({ name: "home" });
        });
});
</script>

