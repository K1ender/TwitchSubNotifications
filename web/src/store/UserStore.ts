import { defineStore } from "pinia";
import { ref } from "vue";
import ky from "ky";

interface User {
  username: string;
}

export const useUserStore = defineStore("user", () => {
  const user = ref<User | null>(null);
  const isAuthorized = ref(false);
  const isLoading = ref(false);

  const setUser = (u: User) => {
    user.value = u;
    isAuthorized.value = !!u;
  };

  const logout = async () => {
    try {
      isLoading.value = true;
      await ky
        .get(import.meta.env.VITE_API_ENDPOINT + "/logout", {
          credentials: "include",
        })
        .json();
      user.value = null;
      isAuthorized.value = false;
    } catch (e) {
      console.error(e);
    } finally {
      isLoading.value = false;
    }
  };

  const getProfile = async () => {
    isLoading.value = true;
    try {
      const res = await ky
        .get(import.meta.env.VITE_API_ENDPOINT + "/profile", {
          credentials: "include",
        })
        .json<{
          success: boolean;
          message: string;
          data: User;
        }>();

      if (res.success) {
        setUser(res.data);
        isAuthorized.value = true;
      }

      return res;
    } catch (e) {
      console.error(e);
    } finally {
      isLoading.value = false;
    }
  };

  return { user, setUser, logout, getProfile, isAuthorized, isLoading };
});
