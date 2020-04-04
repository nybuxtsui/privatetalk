import Vue from 'vue';
import {
  Button,
  Select,
  Row,
  Col,
  Input,
  Table,
  TableColumn,
  Container,
  Header,
  Aside,
  Main,
  Footer,
} from 'element-ui';
import App from './App.vue';

Vue.use(Button);
Vue.use(Select);
Vue.use(Row);
Vue.use(Col);
Vue.use(Input);
Vue.use(Table);
Vue.use(TableColumn);
Vue.use(Container)
Vue.use(Header)
Vue.use(Aside)
Vue.use(Main)
Vue.use(Footer)

Vue.config.productionTip = false

new Vue({
  render: h => h(App),
}).$mount('#app')
