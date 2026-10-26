<script setup lang="ts">
import { ref } from 'vue'
import { onMounted } from 'vue'
import { useUsers } from '@/composables/useUsers'

const newUser = ref({
    name: ''
  })

const selectedUserId = defineModel<number | null>('selectedUserId')
const isAddingUser = defineModel<boolean>('isAddingUser', { default: false })

const { users, error, isLoading, fetchUsers, createUser } = useUsers()

async function handleSubmit() {
    await createUser(newUser.value)    
    newUser.value = { name: '' }
    isAddingUser.value = false
  }

onMounted(fetchUsers)
</script>

<template>
  <div class="user-picker" style="margin-left: 0.5rem; margin-bottom: 1rem;">
    <button @click="isAddingUser = !isAddingUser"> 
        {{ isAddingUser ? 'Cancel' : 'Add User' }} 
    </button>
    <select id="user-select" v-model="selectedUserId" style="margin-left: 1.9rem;">
      <option v-for="user in users" :key="user.id" :value="user.id">
        {{ user.name }}
      </option>
    </select>
    <div v-if="isAddingUser" style="margin-top: 0.5rem; margin-left: -1rem;">
      <form @submit.prevent="handleSubmit">
        <button type="submit" style="margin-left: 1rem;">Add User</button>
        <input v-model="newUser.name" placeholder="Name" required style="margin-left: 1rem;" />
      </form>
    </div>
  </div>
</template>