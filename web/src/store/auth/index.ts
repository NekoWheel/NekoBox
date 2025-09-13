import {defineStore} from "pinia";
import {type AuthState} from "./types";
import type {UserProfile} from "@/api/auth.ts";

const useAuthStore = defineStore('auth', {
    persist: true,

    state: (): AuthState => ({
        isSignedIn: false,
        profile: {} as UserProfile,
        sessionID: '',
    }),

    actions: {
        signIn(profile: UserProfile, sessionID: string) {
            this.isSignedIn = true;
            this.profile = profile;
            this.sessionID = sessionID;
        },
        setProfileName(name: string) {
            this.profile.name = name;
        },
        signOut() {
            this.isSignedIn = false;
            this.profile = {} as UserProfile;
            this.sessionID = '';
        }
    }
})

export default useAuthStore;