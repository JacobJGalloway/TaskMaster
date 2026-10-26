import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { useTasks } from '@/composables/useTasks'
import type { StudyTask } from '@/types/StudyTask'

const sampleTask: StudyTask = {
  id: 1,
  user_id: 1,
  title: 'Read chapter 1',
  duration_minutes: 30,
  period_start: '2026-01-01T09:00',
  period_end: '2026-01-01T09:30',
  completed: false,
  completed_at: null,
}

beforeEach(() => {
  vi.stubGlobal('fetch', vi.fn())
})

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('useTasks', () => {
  it('starts with empty state', () => {
    const { tasks, error, isLoading } = useTasks()
    expect(tasks.value).toEqual([])
    expect(error.value).toBeNull()
    expect(isLoading.value).toBe(false)
  })

  it('fetchTasks hydrates tasks from the API response', async () => {
    vi.mocked(fetch).mockResolvedValue({
      ok: true,
      json: async () => [sampleTask],
    } as Response)

    const { tasks, error, isLoading, fetchTasks } = useTasks()
    const promise = fetchTasks(1)
    expect(isLoading.value).toBe(true)

    await promise

    expect(fetch).toHaveBeenCalledWith('http://localhost:8080/users/1/tasks')
    expect(tasks.value).toEqual([sampleTask])
    expect(error.value).toBeNull()
    expect(isLoading.value).toBe(false)
  })

  it('fetchTasks records an error when the request fails', async () => {
    vi.mocked(fetch).mockResolvedValue({
      ok: false,
      statusText: 'Internal Server Error',
    } as Response)

    const { tasks, error, fetchTasks } = useTasks()
    await fetchTasks(1)

    expect(tasks.value).toEqual([])
    expect(error.value).toContain('Internal Server Error')
  })

  it('createTask posts the new task and refreshes the list', async () => {
    vi.mocked(fetch)
      .mockResolvedValueOnce({ ok: true } as Response) // POST /tasks
      .mockResolvedValueOnce({ ok: true, json: async () => [sampleTask] } as Response) // refresh GET

    const { tasks, createTask } = useTasks()
    await createTask(1, {
      title: 'Read chapter 1',
      duration_minutes: 30,
      period_start: '2026-01-01T09:00',
      period_end: '2026-01-01T09:30',
    })

    expect(fetch).toHaveBeenNthCalledWith(
      1,
      'http://localhost:8080/tasks',
      expect.objectContaining({ method: 'POST' }),
    )
    expect(fetch).toHaveBeenNthCalledWith(2, 'http://localhost:8080/users/1/tasks')
    expect(tasks.value).toEqual([sampleTask])
  })

  it('completeTask posts completion and refreshes the list', async () => {
    vi.mocked(fetch)
      .mockResolvedValueOnce({ ok: true } as Response) // POST /tasks/:id/complete
      .mockResolvedValueOnce({ ok: true, json: async () => [] } as Response) // refresh GET

    const { tasks, completeTask } = useTasks()
    await completeTask(1, 1)

    expect(fetch).toHaveBeenNthCalledWith(
      1,
      'http://localhost:8080/tasks/1/complete',
      expect.objectContaining({ method: 'POST' }),
    )
    expect(tasks.value).toEqual([])
  })
})
