import { ref } from 'vue'
import type { User } from '@/types/User'

export function useUsers() {
  const users = ref<User[]>([])
  const error = ref<string | null>(null)
  const isLoading = ref(false)

    async function fetchUsers() {
        try {
            isLoading.value = true
            const response = await fetch('http://localhost:8080/users')
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.statusText}`)
            }
            users.value = await response.json()
        } catch (err) {
            error.value = err instanceof Error ? err.message : 'An error occurred while fetching users.'
        } finally {
            isLoading.value = false     
        }
    }

    async function createUser(newUser: { name: string }) {
        try {
            isLoading.value = true
            const response = await fetch('http://localhost:8080/users', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(newUser)
            })
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.statusText}`)
            }
            await fetchUsers() // Refresh the user list after creating a new user 
        } catch (err) {
            error.value = err instanceof Error ? err.message : 'An error occurred while creating user.'
        }
        finally {
            isLoading.value = false
        }
    }

    return {   
        users,
        error,
        isLoading,
        fetchUsers,
        createUser
    }
}