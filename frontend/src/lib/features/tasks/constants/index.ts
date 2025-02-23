import type { Category } from "../types";

export const categories: Category[] = [
    {
        value: "focus",
        label: "Focus",
        color: "bg-red-100 dark:bg-red-900/20",
    },
    {
        value: "goals",
        label: "Goals",
        color: "bg-blue-100 dark:bg-blue-900/20",
    },
    {
        value: "fitin",
        label: "Fit In",
        color: "bg-green-100 dark:bg-green-900/20",
    },
    {
        value: "backburner",
        label: "Backburner",
        color: "bg-yellow-100 dark:bg-yellow-900/20",
    },
];

export const categoryTextColors: Record<string, string> = {
    focus: "text-red-700 dark:text-red-300",
    goals: "text-blue-700 dark:text-blue-300",
    fitin: "text-green-700 dark:text-green-300",
    backburner: "text-yellow-700 dark:text-yellow-300",
};