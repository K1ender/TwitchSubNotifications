<script setup lang="ts">
import { motion } from 'motion-v';

const { text } = defineProps<{
    text: string
}>()
const splitText = text.split(" ")
</script>

<template>
    <motion.span v-for="(char, idx) in splitText" :initial="{ opacity: 0, y: -100 }" :animate="{ opacity: 1, y: 0 }"
        :transition="{ delay: idx * 0.2 }" :key="char">
        {{ char + " " }}
    </motion.span>
</template>