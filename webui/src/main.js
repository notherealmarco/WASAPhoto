import {createApp, reactive} from 'vue'
import App from './App.vue'
import router from './router'
import { axios, updateToken as axiosUpdate } from './services/axios.js';
import ErrorMsg from './components/ErrorMsg.vue'
import LoadingSpinner from './components/LoadingSpinner.vue'
import PostCard from './components/PostCard.vue'
import 'bootstrap-icons/font/bootstrap-icons.css'

import './assets/dashboard.css'
import './assets/main.css'

const app = createApp(App)
app.config.globalProperties.$axios = axios;
app.config.globalProperties.$axiosUpdate = axiosUpdate;
app.component("ErrorMsg", ErrorMsg);
app.component("LoadingSpinner", LoadingSpinner);
app.component("PostCard", PostCard);
app.use(router)
app.mount('#app')