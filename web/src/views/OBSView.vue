<script setup lang="ts">
import { useRoute } from 'vue-router';

const router = useRoute();
let websocket = new WebSocket(`ws://${import.meta.env.VITE_WS_ENDPOINT}/obs?id=${router.query.id}`);
websocket.onmessage = (event) => {
    console.log(event.data);
}
</script>

<template>

</template>