export interface Token {
    id: string;
    provider: string;
    account: string;
    access_token: string;
    token_type: string;
    refresh_token?: string;
    expiry?: string;
    scope?: string;
    is_active: boolean;
    created: string;
    user: string;
}

export const DEFAULT_TOKEN_FORM = {
    provider: "",
    account: "",
    access_token: "",
    token_type: "Bearer",
    refresh_token: "",
    expiry: "",
    scope: "",
    is_active: true
};

export const PROVIDERS = [
    { value: "google", label: "Google" },
    { value: "coolify", label: "Coolify" },
    { value: "github", label: "GitHub" },
    { value: "gitlab", label: "GitLab" }
] as const;