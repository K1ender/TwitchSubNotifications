<script setup lang="ts">
import { motion } from 'motion-v';
import { onMounted, ref } from 'vue';
import { useRoute } from 'vue-router';
import { AnimatePresence } from "motion-v"

const router = useRoute();
let websocket = new WebSocket(`ws://${import.meta.env.VITE_WS_ENDPOINT}/obs?id=${router.query.id}`);
const currentEvent = ref<{
    type: string, data: {
        username: string
    }
} | undefined>();

const events = ref<{
    type: string, data: {
        username: string
    }
}[]>([]);

websocket.onmessage = (event) => {
    const json = JSON.parse(event.data);
    if (json.type === "new_subscriber") {
        events.value.push(json);
    }
}

function animateCurrentEvent() {
    if (events.value.length === 0) {
        currentEvent.value = undefined;
        requestAnimationFrame(() => animateCurrentEvent());
        return;
    }

    currentEvent.value = undefined;

    requestAnimationFrame(() => {
        currentEvent.value = events.value.shift();
        setTimeout(animateCurrentEvent, 3000);
    });
}

onMounted(() => {
    animateCurrentEvent();
})
</script>

<template>
    <div class="h-screen bg-transparent w-full flex items-center justify-center">
        <AnimatePresence>
            <motion.div :key="currentEvent.data.username" :initial="{ opacity: 0, scale: 0 }"
                :animate="{ opacity: 1, scale: 1 }" :exit="{ opacity: 0, scale: 0 }" v-if="currentEvent">
                <p class="text-white text-3xl font-bold">
                    {{ currentEvent.data.username }} just subscribed!
                </p>
            </motion.div>
        </AnimatePresence>
    </div>
</template>
