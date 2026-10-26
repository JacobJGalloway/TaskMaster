<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useTasks } from '@/composables/useTasks'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';

const props = defineProps<{ userId: number | null }>()

const newTask = ref({
    title: '',
    duration_minutes: 0,
    period_start: '',
    period_end: ''
  })
const isAddingTask = ref(false)
const { tasks, error, isLoading, fetchTasks, completeTask, createTask } = useTasks()

async function handleSubmit() {
    if (props.userId === null) return
    await createTask(props.userId, newTask.value)
    newTask.value = { title: '', duration_minutes: 0, period_start: '', period_end: '' }
    isAddingTask.value = false
  }

onMounted(() => {
    if (props.userId !== null) {
        fetchTasks(props.userId)
    }
});

watch(
  () => props.userId,
  (newUserId) => {
    if (newUserId !== null) {
      fetchTasks(newUserId)
    }
  }
);

</script>

<template>
  <button @click="isAddingTask = !isAddingTask" style="margin-left: 0.5rem;">
      {{ isAddingTask ? 'Cancel' : 'Add Task' }}
  </button>
  <div v-if="!isAddingTask" class="task-list">
    <ul style="list-style-type:none;">
      <li v-for="task in tasks" :key="task.id">
        <button @click="completeTask(props.userId!, task.id)" size="xl" style="margin-right: 0.5rem; margin-bottom: 0.5rem; padding: 0.25rem; background: none; border: none; cursor: pointer; display: inline-flex; align-items: center; justify-content: center; color: #333;">
          <font-awesome-icon :icon="['far', 'circle-check']" />
        </button>
        {{ task.title }} - {{ task.duration_minutes }} minutes ({{ task.period_start }} to {{ task.period_end }})
      </li>
    </ul>
  </div>  
  <div v-else>
    <form @submit.prevent="handleSubmit" style="margin-left: 0.5rem; margin-top: 0.5rem;">
      <div style="margin-top: 1rem;">
        <div>
          <input v-model="newTask.title" placeholder="Task Title" required style="margin-right: 0.5rem; margin-bottom: 1rem;"/>
          <input v-model.number="newTask.duration_minutes" type="number" placeholder="Duration (minutes)" required />
        </div>
        <div>
          <input v-model="newTask.period_start" type="datetime-local" placeholder="Start Time" required style="margin-right: 0.5rem; margin-bottom: 1rem;"/>
          <input v-model="newTask.period_end" type="datetime-local" placeholder="End Time" required />
        </div>
      </div>
      <button type="submit" style="margin-left: 0rem; margin-top: 0.5rem;">Add Task</button>
      </form>
  </div>

  
  
</template>