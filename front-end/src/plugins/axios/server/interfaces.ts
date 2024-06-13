import type { Server } from "@/components/views/interfaces";

export interface IGetUserByEmailRequest {
    email: string;
}

export interface IUser {
    fullName: string;
    email: string;
    phone?: string;
    avatar?: string;
    role: string;
}

export enum ServerStatus {
    STATUSNONE = 0,
    ON = 1,
    OFF = 2,
}

export interface ICreateServerRequest {
    name: string;
    status: Boolean;
    ip: string;
}

export interface IImportServerRequest {
    listserver: FormData;
}

export interface IServeStatusResopnse {
    online: number;
    offline: number;
}

export enum TypeSort {
    NONE = 0,
    ASC = 1,
    DESC = 2,
}

export interface IServerPagination {
    limit?: number;
    offset?: number;
    status?: string;
    field?: string;
    order?: TypeSort;
}

export interface IServerFilter {
    createdAtFrom?: string;
    createdAtTo?: string;
    updatedAtFrom?: string;
    updatedAtTo?: string;
    status?: ServerStatus;
}

export interface IListServerRequest {
    limit?: number;
    offset?: number;
    status?: string;
    field?: string;
    order?: string;
}

export interface IListServerResponse {
    data: Server[];
    total: number;
}

export interface IUpdateServerRequest {
    id: number;
    name: string;
    status?: ServerStatus;
    ipv4: string;
}

export interface IFileExportServerRequest {
    fileName: string;
}

export interface IPaginationExportRequest {
    pageSize?: number;
    fromPage?: number;
    toPage?: number;
    sort?: TypeSort;
    sortBy?: string;
}

export interface IExportServerRequest {
    limit?: number;
    offset?: number;
    status?: string;
    field?: string;
    order?: string;
}
