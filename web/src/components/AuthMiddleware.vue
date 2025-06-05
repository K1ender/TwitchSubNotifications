<script setup lang="ts">
import { useRouter } from 'vue-router';
import { useUserStore } from '../store/UserStore';
import { onMounted } from 'vue';
const { getProfile } = useUserStore();

const router = useRouter();

onMounted(async () => {
    const user = await getProfile();
    if (!user) {
        return;
    }
    if (user.success) {
        router.push("/dashboard")
    }
})
</script>

<template>
    <slot></slot>
</template>