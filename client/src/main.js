import Vue from 'vue';
import {
  Button,
  Select,
  Row,
  Col,
  Input,
} from 'element-ui';
import App from './App.vue';

Vue.component(Button.name, Button);
Vue.component(Select.name, Select);
Vue.component(Row.name, Row);
Vue.component(Col.name, Col);
Vue.component(Input.name, Input);

Vue.config.productionTip = false

new Vue({
  render: h => h(App),
}).$mount('#app')
