export enum Status {
    ON = "ON",
    OFF = "OFF",
}
export interface Server {
    CreatedAt: string;
    DeletedAt: string;
    ID: number;
    UpdatedAt: string;
    ip: string;
    name: string;
    status: boolean;
}

