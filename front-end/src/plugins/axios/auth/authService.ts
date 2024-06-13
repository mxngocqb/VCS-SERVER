import type { AxiosInstance } from "axios";
import axios from "axios";
import type {
    ILoginRequest,
    ILoginResponse,
    IRefressTokenRequest,
    IRefressTokenResponse,
} from "./interfaces";
import {
    errorHandlerRequest,
    preHandlerResquest,
} from "../interpretors/requestPretorCommon";
import {
    errorHandlerResponse,
    successHandlerResponse,
} from "../interpretors/responsePretorCommon";
import type {
    InterPretorRequest,
    InterPretorResponse,
} from "../interpretors/interfaces";

class AuthService {
    axiosInstance: AxiosInstance;
    constructor(
        baseURL: string,
        interpretorsRequest?: Array<InterPretorRequest>,
        interpretorsResponse?: Array<InterPretorResponse>
    ) {
        this.axiosInstance = axios.create({
            baseURL: baseURL,
            headers: {
                Accept: "application/json",
                "Content-Type": "application/json",
            },
        });
        this.axiosInstance.interceptors.request.use(
            preHandlerResquest,
            errorHandlerRequest
        );
        this.axiosInstance.interceptors.response.use(
            successHandlerResponse,
            errorHandlerResponse
        );
        interpretorsRequest?.forEach((element) => {
            this.axiosInstance.interceptors.request.use(
                element.beforeRequest,
                element.errorHandler
            );
        });
        interpretorsResponse?.forEach((element) => {
            this.axiosInstance.interceptors.response.use(
                element.beforeResponse,
                element.errorHandler
            );
        });
    }

    login(req: ILoginRequest) {
        return this.axiosInstance.post<ILoginResponse>("/login", req);
    }

    refreshToken(req: IRefressTokenRequest) {
        return this.axiosInstance.post<IRefressTokenResponse>(
            "/refresh-token",
            req
        );
    }
}

export const authService = new AuthService(
    `${import.meta.env["VITE_GT_BASE_URL"]}/`
);
