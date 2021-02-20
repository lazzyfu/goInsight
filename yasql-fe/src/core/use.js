import Vue from 'vue'

// base library
import Viser from 'viser-vue'
import VueCropper from 'vue-cropper'

// element
import ElementUI from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'

import Antd from 'ant-design-vue'
import 'ant-design-vue/dist/antd.less'

import '@fortawesome/fontawesome-free/css/all.css'

// codemirror
import VueCodeMirror from 'vue-codemirror'
import 'codemirror/lib/codemirror.css'

// contentmenu
import contentmenu from 'v-contextmenu'
import 'v-contextmenu/dist/index.css'

// jquery
import jQuery from 'jquery'

// echarts
import * as echarts from 'echarts'
Vue.prototype.$echarts = echarts

Vue.use(ElementUI)
Vue.use(Antd)
Vue.use(VueCropper)
Vue.use(VueCodeMirror)
Vue.use(Viser)
Vue.use(contentmenu)

global.$ = jQuery
