import { useMainStore } from "@/stores/mainStore";
import type { AxiosError, AxiosResponse } from "axios";
import { authService } from "../auth/authService";
import { SessionStorageKey } from "@/stores/constants";
import axios from "axios";

export const errorHandlerResponse = async function (error: any) {
    // Bất kì mã trạng thái nào lọt ra ngoài tầm 2xx đều khiến hàm này được trigger\
    // Làm gì đó với lỗi response
    const originalConfig = error.config;
    if (originalConfig?.url !== "/login" && error.response) {
        // Access Token was expired
        if (error.response.status === 401 && !originalConfig?._retry) {
            originalConfig._retry = true;
            await authService
                .refreshToken({
                    refreshToken:
                        sessionStorage.getItem(
                            SessionStorageKey.AUTH_REFRESH_TOKEN
                        ) || "",
                })
                .then((res) => {
                    const { data } = res;
                    sessionStorage.setItem(
                        SessionStorageKey.AUTH_TOKEN,
                        data?.token
                    );
                })
                .catch((err) => {
                    return Promise.reject(error);
                });
            originalConfig.headers.Authorization = `Bearer ${sessionStorage.getItem(
                SessionStorageKey.AUTH_TOKEN
            )}`;

            console.log(originalConfig.headers.Authorization);
            return axios(originalConfig);
        }
    }
    const mainStore = useMainStore();
    mainStore.showNofitication({
        typeNotification: "error",
        title: error.request?.statusText || "",
        content: error.response?.data["message"] || "",
    });
    return Promise.reject(error);
};

export const successHandlerResponse = function (response: AxiosResponse) {
    // Bất kì mã trạng thái nào nằm trong tầm 2xx đều khiến hàm này được trigger
    // Làm gì đó với dữ liệu response
    return response;
};
