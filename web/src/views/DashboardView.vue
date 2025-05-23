<script setup lang="ts">
import Button from '@/components/ui/button/Button.vue';
import Card from '@/components/ui/card/Card.vue';
import CardContent from '@/components/ui/card/CardContent.vue';
import CardHeader from '@/components/ui/card/CardHeader.vue';
import { useUserStore } from '@/store/UserStore';
import { onMounted } from 'vue';
import { useRouter } from 'vue-router';

const store = useUserStore();
const router = useRouter();

onMounted(() => {
    if (!store.isAuthorized && !store.isLoading) {
        router.push("/");
    }
})

async function logout() {
    await store.logout();
    router.push("/");
}
</script>

<template>
    <div class="flex items-center justify-center pt-10">
        <div class="flex flex-col container mx-6">
            <h1 class="text-4xl">{{ store.user?.username }}</h1>
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 my-4">
                <Card class="w-full">
                    <CardHeader>Follower notifications</CardHeader>
                    <CardContent>
                        <Button variant="secondary">Enable</Button>
                    </CardContent>
                </Card>
            </div>
            <Button class="mt-4 cursor-pointer" variant="destructive" @click="logout">Logout</Button>
        </div>
    </div>
</template>