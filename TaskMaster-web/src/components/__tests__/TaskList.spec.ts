import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { ref } from 'vue'
import TaskList from '@/components/TaskList.vue'
import type { StudyTask } from '@/types/StudyTask'

const fetchTasks = vi.fn()
const completeTask = vi.fn()
const createTask = vi.fn()
const tasks = ref<StudyTask[]>([])
const error = ref<string | null>(null)
const isLoading = ref(false)

vi.mock('@/composables/useTasks', () => ({
  useTasks: () => ({ tasks, error, isLoading, fetchTasks, completeTask, createTask }),
}))

const stubs = { FontAwesomeIcon: true }

beforeEach(() => {
  tasks.value = []
  error.value = null
  isLoading.value = false
  fetchTasks.mockClear()
  completeTask.mockClear()
  createTask.mockClear()
})

describe('TaskList', () => {
  it('fetches tasks for the given user on mount', () => {
    mount(TaskList, { props: { userId: 1 }, global: { stubs } })
    expect(fetchTasks).toHaveBeenCalledWith(1)
  })

  it('does not fetch when no user is selected', () => {
    mount(TaskList, { props: { userId: null }, global: { stubs } })
    expect(fetchTasks).not.toHaveBeenCalled()
  })

  it('refetches when the selected user changes', async () => {
    const wrapper = mount(TaskList, { props: { userId: 1 }, global: { stubs } })
    fetchTasks.mockClear()

    await wrapper.setProps({ userId: 2 })

    expect(fetchTasks).toHaveBeenCalledWith(2)
  })

  it('renders each task from composable state', () => {
    tasks.value = [
      {
        id: 1,
        user_id: 1,
        title: 'Read chapter 1',
        duration_minutes: 30,
        period_start: 's',
        period_end: 'e',
        completed: false,
      },
    ]
    const wrapper = mount(TaskList, { props: { userId: 1 }, global: { stubs } })

    expect(wrapper.text()).toContain('Read chapter 1')
    expect(wrapper.text()).toContain('30 minutes')
  })

  it('calls completeTask with the user and task id when the checkmark is clicked', async () => {
    tasks.value = [
      {
        id: 7,
        user_id: 1,
        title: 'Read chapter 1',
        duration_minutes: 30,
        period_start: 's',
        period_end: 'e',
        completed: false,
      },
    ]
    const wrapper = mount(TaskList, { props: { userId: 1 }, global: { stubs } })

    // buttons[0] is the "Add Task" toggle; buttons[1] is the task's complete checkmark.
    await wrapper.findAll('button')[1]!.trigger('click')

    expect(completeTask).toHaveBeenCalledWith(1, 7)
  })

  it('submits a new task through createTask and resets the form', async () => {
    const wrapper = mount(TaskList, { props: { userId: 1 }, global: { stubs } })

    await wrapper.find('button').trigger('click') // reveal the "Add Task" form
    await wrapper.find('input[placeholder="Task Title"]').setValue('Read chapter 2')
    await wrapper.find('input[type="number"]').setValue(45)
    await wrapper.findAll('input[type="datetime-local"]')[0]!.setValue('2026-01-02T09:00')
    await wrapper.findAll('input[type="datetime-local"]')[1]!.setValue('2026-01-02T09:45')

    await wrapper.find('form').trigger('submit.prevent')

    expect(createTask).toHaveBeenCalledWith(1, {
      title: 'Read chapter 2',
      duration_minutes: 45,
      period_start: '2026-01-02T09:00',
      period_end: '2026-01-02T09:45',
    })
  })
})
