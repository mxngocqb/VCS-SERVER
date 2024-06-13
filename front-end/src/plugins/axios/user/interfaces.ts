export interface IGetUserByEmailRequest {
    username: string;
}

export interface IUser {
    id: number;
    fullName: string;
    email: string;
    phone?: string;
    avatar?: string;
    roles: string;
    createdAt?: Date;
    createdBy?: number;
}

export interface IResponseGetListUser {
    users: IUser[];
}

export interface IRequestCreateUser {
    fullName: string;
    email: string;
    password: string;
    phone?: string;
    avatar?: string;
    roles: number;
}

export interface IRequestUpdateUser {
    fullName: string;
    email: string;
    phone?: string;
    avatar?: string;
    roles: number;
}
