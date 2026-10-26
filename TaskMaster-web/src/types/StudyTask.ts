export type StudyTask = {
    id: number;
    user_id: number;
    title: string;
    duration_minutes: number;
    period_start: string;
    period_end: string;
    completed: boolean;
    completed_at?: string | null;
}