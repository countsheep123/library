import Vue from 'vue';
import App from './App.vue';
import router from './router';
import './registerServiceWorker';

import { library } from '@fortawesome/fontawesome-svg-core';
import { fas } from '@fortawesome/free-solid-svg-icons';
import { fab } from '@fortawesome/free-brands-svg-icons';
import { far } from '@fortawesome/free-regular-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import vSelect from 'vue-select';
import 'vue-select/dist/vue-select.css';

Vue.component('v-select', vSelect);

library.add(fas, far, fab);
Vue.component("font-awesome-icon", FontAwesomeIcon);

import VueQuagga from "vue-quaggajs";
Vue.use(VueQuagga);

Vue.config.productionTip = false;

new Vue({
	router,
	render: h => h(App)
}).$mount('#app');
