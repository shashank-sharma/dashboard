export interface Cell {
    id: string;
    content: string;
    output: string;
    type: 'code' | 'markdown';
    language: string;
}

export interface Notebook {
    id: string;
    name: string;
    version: string;
    cells: Cell[];
    created: string;
    updated: string;
    user: string;
}

export interface CreateNotebookData {
    name: string;
    version: string;
    cells: Cell[];
}