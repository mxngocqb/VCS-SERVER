import { useMainStore } from "@/stores/mainStore";
import type { InternalAxiosRequestConfig } from "axios";

export const errorHandlerRequest = function (error: any) {
    // Làm gì đó với lỗi request
    const mainStore = useMainStore();
    mainStore.showNofitication({
        typeNotification: "error",
        title: error.request?.statusText || "",
        content: error.response?.data["message"] || "",
    });
    return Promise.reject(error);
};

export const preHandlerResquest = function (
    config: InternalAxiosRequestConfig<any>
) {
    // Làm gì đó trước khi request dược gửi đi
    return config;
};
