export interface ILoginRequest {
    username: string;
    password: string;
}

export interface ILoginResponse {
    token: string;
    expireTime: string;
    typeToken: string;
    refreshToken: string;
}

export interface IRefressTokenRequest {
    refreshToken: string;
}

export interface IRefressTokenResponse {
    token: string;
}
