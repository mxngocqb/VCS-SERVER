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
import type { ISendReportRequest } from "./interfaces";

class MailService {
    axiosInstance: AxiosInstance;

    constructor(
        baseURL: string,
        interpretorsRequest?: Array<InterPretorRequest>,
        interpretorsResponse?: Array<InterPretorResponse>
    ) {
        // Ensure baseURL does not have a trailing slash
        if (baseURL.endsWith('/')) {
            baseURL = baseURL.slice(0, -1);
        }

        this.axiosInstance = axios.create({
            baseURL: baseURL,
            headers: {
                Accept: "application/json",
            },
        });

        // Add interceptors if provided
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

    sendReport(data: ISendReportRequest) {
        // Format dates to YYYY-MM-DD
        const start = encodeURIComponent(data.from.split('T')[0]);
        const end = encodeURIComponent(data.to.split('T')[0]);

        // Construct query parameters
        const queryParams = `report?start=${start}&end=${end}&mail=${encodeURIComponent(data.email)}`;

        console.log(this.axiosInstance.defaults.baseURL + queryParams); // Log the final URL for debugging purposes

        this.axiosInstance.get(queryParams)
            .then(response => {
                // Handle response
                console.log(response.data);
            })
            .catch(error => {
                // Handle error
                console.error(error);
            });
    }
}

// Ensure no trailing slash in the base URL during instantiation
export const mailService = new MailService(
    `${import.meta.env["VITE_GT_BASE_URL"]}/servers`,
    [authInterpretor],
    []
);


