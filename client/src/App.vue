<template>
  <header>Chat App</header>
  <main>
    <template v-if="isConnected">
      <Chat
        v-model="newMessage"
        :messages="messages"
        :id="id"
        @send-message="sendMessage"
      />
    </template>
    <template v-else>
      <SignIn v-model="name" @connect="connect" />
    </template>
  </main>
</template>

<script setup lang="ts">
import {ref, onBeforeUnmount} from "vue";
import SignIn from "./components/SignIn.vue";
import Chat from "./components/Chat.vue";
import type {MessageWrapper} from "@/types/chat";

const name = ref<string>("User");
const id = ref<string>("");
const messages = ref<MessageWrapper[]>([]);
const isConnected = ref(false);
const newMessage = ref("");
const socket = ref<WebSocket | null>(null);

const sendMessage = () => {
  if (socket.value && newMessage.value.trim()) {
    const data = {
      name: name.value,
      message: newMessage.value,
    };
    socket.value.send(JSON.stringify(data));
    newMessage.value = "";
  }
};

const connect = () => {
  socket.value = new WebSocket("ws://localhost:8080/ws");

  socket.value.onopen = () => {
    console.log("Connected");
    isConnected.value = true;
  };

  socket.value.onmessage = (event) => {
    try {
      const parsedData = JSON.parse(event.data);
      if (parsedData.senderId == "server") {
        id.value = parsedData.content;
        return;
      }
      if (typeof parsedData.content === "string") {
        parsedData.content = JSON.parse(parsedData.content);
      }
      messages.value.push(parsedData);
    } catch (error) {
      console.error("Error parsing WebSocket message:", error);
    }
  };

  socket.value.onclose = () => {
    console.log("Disconnected from server");
    isConnected.value = false;
  };

  socket.value.onerror = (error) => {
    console.error("WebSocket error:", error);
  };
};

onBeforeUnmount(() => {
  if (socket.value) {
    socket.value.close();
  }
});
</script>
