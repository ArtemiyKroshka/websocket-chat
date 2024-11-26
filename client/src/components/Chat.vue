<template>
  <div class="chat">
    <div class="chat-window">
      <div
        class="chat-window__block"
        :class="{'from-me': props.id === el.senderId}"
        v-for="(el, idx) of props.messages"
        :key="idx"
      >
        <span class="name">
          {{ el.content.name }}
        </span>
        <span class="message">
          {{ el.content.message }}
        </span>
      </div>
    </div>
    <form @submit.prevent="emit('sendMessage')" class="chat-input">
      <input type="text" v-model="message" />
      <button type="submit">Send</button>
    </form>
  </div>
</template>
<script lang="ts" setup>
import {ref, watch, defineProps, defineEmits} from "vue";
import type {MessageWrapper} from "@/types/chat";

type Props = {
  modelValue: string;
  messages: MessageWrapper[];
  id: string;
};

const props = defineProps<Props>();

const emit = defineEmits(["update:modelValue", "sendMessage"]);

const message = ref<string>(props.modelValue);

watch(message, (newValue) => {
  emit("update:modelValue", newValue);
});

watch(
  () => props.modelValue,
  () => {
    message.value = props.modelValue;
  }
);
</script>

<style lang="scss" scoped>
.chat {
  display: flex;
  flex-direction: column;
  width: 100%;
  justify-content: center;
  gap: 12px 0;

  &-window {
    border-radius: 10px;
    min-height: 300px;
    max-height: 350px;
    overflow-y: auto;
    background-color: var(--vt-c-text-dark-2);
    display: flex;
    flex-direction: column;
    gap: 12px 0;
    padding: 24px;

    &__block {
      display: flex;
      flex-direction: column;
      justify-content: start;
      width: fit-content;
      border-radius: 10px;
      padding: 8px 12px;
      background-color: var(--vt-c-black-soft);
      color: var(--vt-c-white-mute);
      max-width: 40%;

      &.from-me {
        margin-left: auto;
      }

      .name {
        font-size: 12px;
        color: var(--color-text);
        font-weight: 600;
      }

      .message {
        word-break: break-all;
      }
    }
  }

  &-input {
    display: flex;
    align-items: center;
    gap: 0 12px;

    input {
      text-align: left;
      padding: 0 8px;
      flex: 1;
    }

    button {
      padding-left: 24px;
      padding-right: 24px;
      font-size: 16px;
      cursor: pointer;

      &:hover,
      &:focus {
        transition: 0.5 all ease;
        background-color: var(--vt-c-black-mute);
      }
    }
  }
}
</style>
