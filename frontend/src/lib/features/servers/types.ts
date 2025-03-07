export interface Server {
    id: string;
    name: string;
    provider: string;
    ip: string;
    is_active: boolean;
    is_reachable: boolean;
    port: number;
    username: string;
    security_key: string;
    ssh_enabled: boolean;
    created: string;
    updated: string;
    user: string;
}

export const DEFAULT_SERVER_FORM = {
    name: "",
    provider: "",
    ip: "",
    port: 22,
    username: "",
    security_key: "",
    ssh_enabled: false,
    is_active: true,
};