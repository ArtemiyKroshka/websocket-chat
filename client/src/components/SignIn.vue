<template>
  <form class="form" @submit.prevent="emit('connect')">
    <label for="name">Enter your name:</label>
    <input name="name" type="text" v-model="localName" />
    <button type="submit">Submit</button>
  </form>
</template>

<script setup lang="ts">
import {ref, watch, defineProps, defineEmits} from "vue";

defineProps({
  modelValue: {
    type: String,
    required: true,
  },
});

const emit = defineEmits(["update:modelValue", "connect"]);

const localName = ref<string>("User");

watch(localName, (newValue) => {
  emit("update:modelValue", newValue);
});
</script>

<style lang="scss" scoped>
.form {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px 24px;
  flex-wrap: wrap;
  max-width: 720px;

  button,
  input {
    flex: 0 0 100%;
  }
}
</style>
