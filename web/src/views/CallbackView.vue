<script setup lang="ts">
import { useUserStore } from '@/store/UserStore';
import ky from 'ky';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
const err = ref('Loading...');

const router = useRouter()

const { getProfile } = useUserStore();

onMounted(async () => {
    const url = new URL(window.location.href);
    const code = url.searchParams.get('code');
    const scope = url.searchParams.get('scope');
    const state = url.searchParams.get('state');
    if (code && scope && state) {
        const url = new URL(import.meta.env.VITE_API_ENDPOINT + "/callback");
        url.searchParams.set('code', code);
        url.searchParams.set('scope', scope);
        url.searchParams.set('state', state);
        const res = await ky.get(url.toString(), { credentials: 'include' });
        await getProfile();
        if (res.ok) {
            router.push("/dashboard")
        } else {
            err.value = await res.json();
        }
    }

})
</script>

<template>
    <p>{{ err }}</p>
</template>