<script setup lang="ts">
import Button from '@/components/ui/button/Button.vue'

async function login() {
    const url = new URL(import.meta.env.VITE_API_ENDPOINT + "/login");
    window.location.href = url.toString();
}
</script>

<template>
    <div class="h-screen w-screen flex items-center justify-center">
        <Button @click="login">Login with Twitch</Button>
    </div>
</template>
