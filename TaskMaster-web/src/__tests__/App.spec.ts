import { describe, it, expect } from 'vitest'
import { nextTick } from 'vue'
import { shallowMount } from '@vue/test-utils'
import App from '../App.vue'
import UserPicker from '@/components/UserPicker.vue'
import TaskList from '@/components/TaskList.vue'

describe('App', () => {
  it('renders the heading', () => {
    const wrapper = shallowMount(App)
    expect(wrapper.text()).toContain('Task Master')
  })

  it('renders TaskList with no user selected initially', () => {
    const wrapper = shallowMount(App)
    const taskList = wrapper.findComponent(TaskList)
    expect(taskList.exists()).toBe(true)
    expect(taskList.props('userId')).toBeNull()
  })

  it('passes the selected user id down to TaskList once chosen', async () => {
    const wrapper = shallowMount(App)
    const picker = wrapper.findComponent(UserPicker)

    picker.vm.$emit('update:selectedUserId', 5)
    await nextTick()

    expect(wrapper.findComponent(TaskList).props('userId')).toBe(5)
  })

  it('hides TaskList while a new user is being added', async () => {
    const wrapper = shallowMount(App)
    const picker = wrapper.findComponent(UserPicker)

    picker.vm.$emit('update:isAddingUser', true)
    await nextTick()

    expect(wrapper.findComponent(TaskList).exists()).toBe(false)
  })
})
