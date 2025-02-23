export interface Server {
    id: string;
    name: string;
    provider: string;
    url: string;
    token: string;
    is_active: boolean;
    created: string;
    updated: string;
    user: string;
}

export const DEFAULT_SERVER_FORM = {
    name: "",
    provider: "",
    url: "",
    token: "",
    is_active: true
};