import { createApp } from 'vue'
import App from './App.vue'

import Toast from "vue-toastification";

// Import the CSS or use your own!
import "vue-toastification/dist/index.css";

import router from './router'

const app = createApp(App).use(router);

app.use(Toast, {});

app.mount("#app");
