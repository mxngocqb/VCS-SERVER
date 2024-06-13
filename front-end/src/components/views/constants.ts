import type {
    TypeSort,
    IListServerRequest,
    IServerPagination
} from "@/plugins/axios/server/interfaces";

export const DefaultQuery: IListServerRequest = {
    limit: 10,
    offset: 0,
    order: "asc",
    field: "created_at",
};

