import {createPinia} from 'pinia';
import useAppStore from './app'
import useAuthStore from "./auth";
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'

const pinia = createPinia();
pinia.use(piniaPluginPersistedstate);

export {useAppStore, useAuthStore}
export default pinia;