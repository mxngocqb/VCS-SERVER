import type { AxiosInstance } from "axios";
import axios from "axios";
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
import { authInterpretor } from "../interpretors/authorizeInterpretor";
import type { Server } from "@/components/views/interfaces";
import qs from "qs";

import type {
    ICreateServerRequest,
    IExportServerRequest,
    IImportServerRequest,
    IListServerRequest,
    IListServerResponse,
    IUpdateServerRequest,
} from "./interfaces";

class ServerService {
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
                // "Access-Control-Allow-Origin": "*",
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
    getListServer(req: IListServerRequest) {
        return this.axiosInstance.get<IListServerResponse>("", {
            params: req,
            paramsSerializer: (param) =>
                qs.stringify(param, {
                    indices: false,
                    encode: false,
                    allowDots: true,
                }),
        });
    }
    getListServerStatus() {
        return this.axiosInstance.get("/status");
    }
    createServer(data: ICreateServerRequest) {
        return this.axiosInstance.post("", data);
    }
    deleteServer(id: number) {
        return this.axiosInstance.delete(`/${id}`);
    }
    updateServer(server: IUpdateServerRequest) {
        return this.axiosInstance.patch(`/${server.id}`, server);
    }

    importServer(file: File) {
        let formData = new FormData();
        formData.append("listserver", file);

        return this.axiosInstance.post("/import", formData, {
            headers: {
                "Content-Type": "multipart/form-data",
            },
        });
    }

    exportServer(req: IExportServerRequest) {
        if (req.limit === undefined) req.limit = 10;
        if (req.offset === undefined) req.offset = 0;
        if (req.status === undefined) req.status = "true";
        if (req.field === undefined) req.field = "id";
        if (req.order === undefined) req.order = "asc";

        console.log(req);
    
        this.axiosInstance.get(`/export?limit=${req.limit}&offset=${req.offset}&status=${req.status}&field=${req.field}&order=${req.order}`, { responseType: 'blob' })
            .then((response) => {
                const url = window.URL.createObjectURL(new Blob([response.data]));
                const link = document.createElement("a");
                link.href = url;
                link.setAttribute("download", "export.xlsx");
                document.body.appendChild(link);
                link.click();
                // Clean up after download
                window.URL.revokeObjectURL(url);
            })
            .catch((error) => {
                console.error('Error downloading file:', error);
            });

    }
}

export const serverService = new ServerService(
    `${import.meta.env["VITE_GT_BASE_URL"]}/servers`,
    [authInterpretor],
    []
);
