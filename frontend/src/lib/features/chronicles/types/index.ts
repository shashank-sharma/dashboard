export interface JournalEntry {
    id?: string;
    user: string;
    title: string;
    content: string;
    date: string;
    mood: 'happy' | 'neutral' | 'sad' | 'excited' | 'anxious' | 'peaceful';
    tags: string;
}

export interface ChronicleState {
    currentEntry: JournalEntry | null;
    isLoading: boolean;
    error: string | null;
}