<script setup lang="ts">
import { ref } from 'vue';
import UserPicker from '@/components/UserPicker.vue';
import TaskList from '@/components/TaskList.vue';

const selectedUserId = ref<number | null>(null);
const isAddingUser = ref(false);
</script>

<template>
  <h1 style="margin-left: 0.5rem;">Task Master</h1>
  <UserPicker v-model:selectedUserId="selectedUserId" v-model:isAddingUser="isAddingUser" />
  <TaskList v-if="!isAddingUser" :userId="selectedUserId" />
</template>

<style scoped></style>
