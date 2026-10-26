import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { ref } from 'vue'
import UserPicker from '@/components/UserPicker.vue'
import type { User } from '@/types/User'

const fetchUsers = vi.fn()
const createUser = vi.fn()
const users = ref<User[]>([])
const error = ref<string | null>(null)
const isLoading = ref(false)

vi.mock('@/composables/useUsers', () => ({
  useUsers: () => ({ users, error, isLoading, fetchUsers, createUser }),
}))

beforeEach(() => {
  users.value = []
  error.value = null
  isLoading.value = false
  fetchUsers.mockClear()
  createUser.mockClear()
})

describe('UserPicker', () => {
  it('fetches users on mount', () => {
    mount(UserPicker)
    expect(fetchUsers).toHaveBeenCalledTimes(1)
  })

  it('renders an option per fetched user', () => {
    users.value = [
      { id: 1, name: 'Ada' },
      { id: 2, name: 'Grace' },
    ]
    const wrapper = mount(UserPicker)

    const options = wrapper.findAll('option')
    expect(options).toHaveLength(2)
    expect(options[0]!.text()).toBe('Ada')
    expect(options[1]!.text()).toBe('Grace')
  })

  it('emits the selected user id via v-model', async () => {
    users.value = [{ id: 1, name: 'Ada' }]
    const wrapper = mount(UserPicker, { props: { selectedUserId: null } })

    await wrapper.find('select').setValue(1)

    expect(wrapper.emitted('update:selectedUserId')?.[0]).toEqual([1])
  })

  it('toggles the add-user form and emits isAddingUser', async () => {
    const wrapper = mount(UserPicker, { props: { isAddingUser: false } })

    await wrapper.find('button').trigger('click')

    expect(wrapper.emitted('update:isAddingUser')?.[0]).toEqual([true])
  })

  it('submits a new user through createUser and resets the form', async () => {
    const wrapper = mount(UserPicker, { props: { isAddingUser: true } })

    await wrapper.find('input[placeholder="Name"]').setValue('Ada')
    await wrapper.find('form').trigger('submit.prevent')

    expect(createUser).toHaveBeenCalledWith({ name: 'Ada' })
  })
})
