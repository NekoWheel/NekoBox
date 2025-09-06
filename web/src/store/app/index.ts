import {defineStore} from "pinia";
import {type AppState} from "./types";

const useAppStore = defineStore('app', {
    persist: true,

    state: (): AppState => ({
        theme: '',
    }),

    actions: {
        setTheme(theme: string) {
            this.theme = theme;
        },
    }
})

export default useAppStore;