import { ref } from "vue"
import type { StudyTask } from "@/types/StudyTask"

export function useTasks() {
  const tasks = ref<StudyTask[]>([])
  const error = ref<string | null>(null)
  const isLoading = ref(false)

    async function fetchTasks(userId: number) {
        try {
            isLoading.value = true
            const response = await fetch(`http://localhost:8080/users/${userId}/tasks`)
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.statusText}`)
            }
            tasks.value = await response.json()
        } catch (err) {
            error.value = err instanceof Error ? err.message : 'An error occurred while fetching tasks.'
        } finally {
            isLoading.value = false
        }
    }

    async function createTask(userId: number, newTask: Omit<StudyTask, 'id' | 'user_id' | 'completed' | 'completed_at'>) {
        try {
            isLoading.value = true
            const sendTask = {
                ...newTask,
                user_id: userId,
                completed: false,
                completed_at: null
            }
            const response = await fetch(`http://localhost:8080/tasks`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(sendTask)
            })
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.statusText}`)
            }
             await fetchTasks(userId) // Refresh the task list after creating a new task
        } catch (err) {
            error.value = err instanceof Error ? err.message : 'An error occurred while creating a task.'
        } finally {
            isLoading.value = false
        }
    }

    async function completeTask(userId: number, taskId: number) {
        try {
            isLoading.value = true
            const response = await fetch(`http://localhost:8080/tasks/${taskId}/complete`, {
                method: 'POST'
            })
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.statusText}`)
            }
            await fetchTasks(userId) // Refresh the task list after completing a task
        } catch (err) {
            error.value = err instanceof Error ? err.message : 'An error occurred while completing the task.'
        } finally {
            isLoading.value = false
        }
    }

    return {
        tasks,
        error,
        isLoading,
        fetchTasks,
        createTask,
        completeTask
    }
}