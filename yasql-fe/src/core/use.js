import Vue from 'vue'

// base library
import Viser from 'viser-vue'
import VueCropper from 'vue-cropper'

// element
import ElementUI from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'

import Antd from 'ant-design-vue'
import 'ant-design-vue/dist/antd.less'


import VueCodeMirror from 'vue-codemirror'
import 'codemirror/lib/codemirror.css'

Vue.use(ElementUI);
Vue.use(Antd)
Vue.use(VueCropper)
Vue.use(VueCodeMirror)
Vue.use(Viser)


process.env.NODE_ENV !== 'production'
