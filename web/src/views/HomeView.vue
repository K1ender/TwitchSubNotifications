<script setup lang="ts">
import AnimatedText from '@/components/AnimatedText.vue';
import Button from '@/components/ui/button/Button.vue';
import Card from '@/components/ui/card/Card.vue';
import CardContent from '@/components/ui/card/CardContent.vue';
import CardDescription from '@/components/ui/card/CardDescription.vue';
import CardHeader from '@/components/ui/card/CardHeader.vue';
import CardTitle from '@/components/ui/card/CardTitle.vue';
import { Twitch } from 'lucide-vue-next';
import { motion } from 'motion-v';

async function login() {
    const url = new URL(import.meta.env.VITE_API_ENDPOINT + "/login");
    window.location.href = url.toString();
}
</script>

<template>
    <div
        class="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900 flex items-center justify-center p-4 relative overflow-hidden">
        <div class="absolute inset-0 overflow-hidden">
            <div class="absolute -top-40 -right-40 w-80 h-80 bg-purple-500/20 rounded-full blur-3xl animate-pulse">
            </div>
            <div
                class="absolute -bottom-40 -left-40 w-80 h-80 bg-blue-500/20 rounded-full blur-3xl animate-pulse delay-1000">
            </div>
            <div
                class="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 w-96 h-96 bg-violet-500/10 rounded-full blur-3xl animate-pulse delay-500">
            </div>
        </div>

        <Card class="w-full max-w-lg bg-slate-800/80 backdrop-blur-xl border-slate-700/50 shadow-2xl relative z-10">
            <CardHeader class="text-center pb-4">
                <div class="flex justify-center mb-6">
                    <div class="relative">
                        <div class="absolute inset-0 bg-purple-500/30 rounded-full blur-xl animate-pulse"></div>
                        <motion.div :initial="{ scale: 0 }" :animate="{ scale: 1 }"
                            :transition="{ duration: 1, type: 'spring' }"
                            class="relative p-4 bg-gradient-to-br from-purple-600 to-blue-600 rounded-full">
                            <Twitch class="h-12 w-12 text-white" />
                        </motion.div>
                    </div>
                </div>
                <CardTitle
                    class="text-3xl font-bold bg-gradient-to-r from-purple-400 to-blue-400 bg-clip-text text-transparent mb-2">
                    Twitch Notifications
                </CardTitle>
                <CardDescription class="text-slate-300 text-lg">
                    <AnimatedText text="Connect with your Twitch account to get started">
                    </AnimatedText>
                </CardDescription>
            </CardHeader>

            <CardContent class="flex justify-center">
                <Button class="w-full bg-gradient-to-r from-purple-600 to-blue-600 hover:from-purple-700 hover:to-blue-700 text-white border-0 shadow-lg transition-all duration-300 flex
                    items-center justify-center h-12 cursor-pointer" @click="login" size="lg">
                    <Twitch class="mr-3 h-5 w-5" />
                    <span>Login with Twitch</span>
                </Button>
            </CardContent>
        </Card>
    </div>
</template>
