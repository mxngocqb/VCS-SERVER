import type { AxiosResponse, InternalAxiosRequestConfig } from "axios";

export class InterPretorRequest {
    beforeRequest(config: InternalAxiosRequestConfig<any>) {
        return config;
    }
    errorHandler(err: any) {}
}

export class InterPretorResponse {
    beforeResponse(response: AxiosResponse<any>) {
        return response;
    }
    errorHandler(err: any) {}
}
