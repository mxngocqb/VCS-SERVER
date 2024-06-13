import type { InternalAxiosRequestConfig } from "axios";
import { InterPretorRequest } from "./interfaces";
import { SessionStorageKey } from "@/stores/constants";

class AuthInterpretor extends InterPretorRequest {
    beforeRequest(
        config: InternalAxiosRequestConfig<any>
    ): InternalAxiosRequestConfig<any> {
        config.headers.Authorization = `Bearer ${sessionStorage.getItem(
            SessionStorageKey.AUTH_TOKEN
        )}`;

        console.log(SessionStorageKey.AUTH_TOKEN);
        
        return config;
    }
}

export const authInterpretor = new AuthInterpretor();
