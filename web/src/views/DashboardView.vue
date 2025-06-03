<script setup lang="ts">
import Button from '@/components/ui/button/Button.vue';
import Card from '@/components/ui/card/Card.vue';
import CardContent from '@/components/ui/card/CardContent.vue';
import CardDescription from '@/components/ui/card/CardDescription.vue';
import CardHeader from '@/components/ui/card/CardHeader.vue';
import CardTitle from '@/components/ui/card/CardTitle.vue';
import { useUserStore } from '@/store/UserStore';
import ky from 'ky';
import { Users, TrendingUp, UserPlus, LogOut, Image } from 'lucide-vue-next'
import { AnimatePresence, motion } from 'motion-v';
import { computed, onMounted, ref } from 'vue';
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
    await router.push("/");
}

const totalFollowers = ref(0);
const newFollowersThisWeek = ref(0);
const followingBack = ref(0);

const followers = ref<{
    id: number,
    displayName: string,
    avatar: string,
    followedAt: string
}[]>([]);
const limit = ref(10);
const offset = ref(0);

const isSubscribed = ref(false);

const subscriptions = ref<{
    id: string,
    type: string,
    condition: {
        broadcaster_user_id?: string,
        moderator_user_id?: string,
        broadcaster_id?: string,
        user_id?: string,
    }
}[]>([]);

async function getSubscriptions() {
    let res = await ky.get(import.meta.env.VITE_API_ENDPOINT + `/subscribed`, {
        credentials: 'include'
    })
    if (!res.ok) {
        return;
    }
    const data = await res.json<
        {
            success: boolean, data: {
                id: string,
                type: string,
                condition: {
                    broadcaster_user_id?: string,
                    moderator_user_id?: string,
                    broadcaster_id?: string,
                    user_id?: string,
                }
            }[],
            message: string
        }>();
    if (data.data.length > 0) {
        isSubscribed.value = true;
    }
    subscriptions.value = data.data;
}

const follow_sub_id = computed(() => subscriptions.value.find(sub => {
    return sub.type === 'channel.follow'
})?.id);

async function getLatestFollowers() {
    let res = await ky.get(import.meta.env.VITE_API_ENDPOINT + `/followers?limit=${limit.value}&offset=${offset.value}`, {
        credentials: 'include'
    });
    if (res.ok) {
        const data = await res.json<
            {
                success: boolean,
                data: { id: number, displayName: string, avatar: string, followedAt: string }[],
                message: string
            }
        >();
        followers.value = data.data;
    }
}

async function subscribe() {
    if (!store.user) {
        return;
    }

    let res = await ky.post(import.meta.env.VITE_API_ENDPOINT + `/subscribe/${store.user.id}`, {
        credentials: 'include'
    });
    if (res.ok) {
        isSubscribed.value = true;
    }
    await getSubscriptions();
}

async function unsubscribe() {
    console.log(follow_sub_id.value);
    let res = await ky.post(import.meta.env.VITE_API_ENDPOINT + `/unsubscribe/${follow_sub_id.value}`, {
        credentials: 'include'
    });
    if (res.ok) {
        isSubscribed.value = false;
    }
    await getSubscriptions();
}

onMounted(async () => {
    await getSubscriptions();
    await getLatestFollowers();
})
</script>

<template>
    <div class="min-h-screen bg-background p-4 md:p-6 lg:p-8 max-w-7xl mx-auto space-y-6">
        <div class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between ">
            <div>
                <h1 class="text-3xl font-bold tracking-tight">
                    Followers
                </h1>
                <p class="text-muted-foreground">
                    Manage and view your follower community
                </p>
            </div>
            <div>
                <Button @click="logout" variant="outline" class="flex items-center gap-2">
                    <LogOut class="h-4 w-4" />
                    Logout
                </Button>
            </div>
        </div>

        <div>
            <Button @click="subscribe" v-if="!isSubscribed" variant="default">
                Subscribe to follows
            </Button>
            <Button @click="unsubscribe" v-else variant="destructive">
                Unsubscribe to follows
            </Button>
        </div>
        <div class="grid gap-4 md:grid-cols-3">
            <div>
                <Card class="hover:shadow-lg transition-shadow duration-300">
                    <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle class="text-sm font-medium">Total Followers</CardTitle>
                        <motion.div :while-hover="{ scale: 1.1, rotate: 5 }"
                            :transition="{ type: 'spring', stiffness: 400, damping: 10 }">
                            <Users class="h-4 w-4 text-muted-foreground" />
                        </motion.div>
                    </CardHeader>
                    <CardContent>
                        <div class="text-2xl font-bold">
                            {{ totalFollowers.toLocaleString() }}
                        </div>
                    </CardContent>
                </Card>
            </div>

            <div>
                <Card class="hover:shadow-lg transition-shadow duration-300">
                    <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle class="text-sm font-medium">New This Week</CardTitle>
                        <motion.div :while-hover="{ scale: 1.1, rotate: 5 }"
                            :transition="{ type: 'spring', stiffness: 400, damping: 10 }">
                            <TrendingUp class="h-4 w-4 text-muted-foreground" />
                        </motion.div>
                    </CardHeader>
                    <CardContent>
                        <div class="text-2xl font-bold">
                            {{ newFollowersThisWeek }}
                        </div>
                    </CardContent>
                </Card>
            </div>

            <div>
                <Card class="hover:shadow-lg transition-shadow duration-300">
                    <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle class="text-sm font-medium">Following Back</CardTitle>
                        <motion.div :while-hover="{ scale: 1.1, rotate: 5 }"
                            :transition="{ type: 'spring', stiffness: 400, damping: 10 }">
                            <UserPlus class="h-4 w-4 text-muted-foreground" />
                        </motion.div>
                    </CardHeader>
                    <CardContent>
                        <div class="text-2xl font-bold">
                            {{ followingBack }}
                        </div>
                    </CardContent>
                </Card>
            </div>
        </div>
        <Card>
            <CardHeader>
                <CardTitle>Follower List</CardTitle>
                <CardDescription>View and manage your followers</CardDescription>
            </CardHeader>
            <CardContent>
                <motion.div class="mt-6 space-y-4" layout>
                    <AnimatePresence mode="popLayout">
                        <motion.div v-for="(follower, index) in followers" layout :key="follower.id"
                            class="flex items-center justify-between p-4 border rounded-lg hover:bg-muted/50 transition-colors cursor-pointer">
                            <div className="flex items-center gap-4">
                                <div>
                                    <img :src="follower.avatar || 'https://placehold.co/40x40'"
                                        :alt="follower.displayName" width="40" height="40" class="rounded-full" />
                                </div>
                                <div class="flex-1 min-w-0">
                                    <div class="flex items-center gap-2">
                                        <p class="font-medium truncate">{{ follower.displayName }}</p>
                                    </div>
                                </div>
                            </div>
                            <div class="flex items-center gap-2">
                                <p class="text-xs text-muted-foreground">
                                    {{ follower.followedAt }}
                                </p>
                                <motion.div :while-hover="{ scale: 1.1, rotate: 5 }"
                                    :transition="{ type: 'spring', stiffness: 400, damping: 10 }">
                                    <Image class="h-4 w-4 text-muted-foreground" />
                                </motion.div>
                            </div>
                        </motion.div>
                    </AnimatePresence>
                </motion.div>
            </CardContent>
        </Card>
    </div>
</template>