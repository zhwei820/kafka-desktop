import { createApp } from "vue";
import App from "./App.vue";
import "./style.css";
import "element-plus/dist/index.css";

// 仅用作示例

const app = createApp(App);

import ElementPlus from "element-plus";
app.use(ElementPlus);

app.mount("#app");
