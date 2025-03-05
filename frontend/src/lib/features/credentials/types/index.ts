import type { Token } from "../../tokens/types";

export interface DeveloperToken {
    id: string;
    name: string;
    token: string;
    is_active: boolean;
    created: string;
    expires?: string;
    user: string;
    environment: string;
}

export interface ApiKey {
    id: string;
    service: string;
    key: string;
    secret?: string;
    is_active: boolean;
    created: string;
    expires?: string;
    user: string;
}

export interface CredentialsStats {
    totalTokens: number;
    totalDeveloperTokens: number;
    totalApiKeys: number;
}

export interface CredentialsState {
    stats: CredentialsStats;
    isLoading: boolean;
}

export const DEFAULT_DEV_TOKEN_FORM = {
    name: "",
    environment: "development",
    is_active: true,
};

export const DEFAULT_API_KEY_FORM = {
    service: "",
    key: "",
    secret: "",
    is_active: true,
};

export const ENVIRONMENTS = [
    { value: "development", label: "Development" },
    { value: "staging", label: "Staging" },
    { value: "production", label: "Production" },
] as const;

export const SERVICES = [
    { value: "github", label: "GitHub" },
    { value: "gitlab", label: "GitLab" },
    { value: "coolify", label: "Coolify" },
    { value: "aws", label: "AWS" },
    { value: "gcp", label: "Google Cloud" },
    { value: "azure", label: "Azure" },
    { value: "openai", label: "OpenAI" },
    { value: "other", label: "Other" },
] as const; 