import Vue from 'vue'

// base library
import Antd from 'ant-design-vue/es'
import Viser from 'viser-vue'
import VueCropper from 'vue-cropper'
import 'ant-design-vue/dist/antd.less'

// ext library
import VueClipboard from 'vue-clipboard2'
import MultiTab from '@/components/MultiTab'
import PageLoading from '@/components/PageLoading'
import PermissionHelper from '@/core/permission/permission'

// jquery
import jQuery from 'jquery'
// 需要定义window.jQuery，否则会抛出Cannot set properties of undefined (setting 'BootstrapTable')
window.jQuery = jQuery
global.$ = jQuery

import './directives/action'

// codemirror
import VueCodeMirror from 'vue-codemirror'
import 'codemirror/lib/codemirror.css'

// contentmenu
import contentmenu from 'v-contextmenu'
import 'v-contextmenu/dist/index.css'

VueClipboard.config.autoSetContainer = true

Vue.use(Antd)
Vue.use(Viser)
Vue.use(MultiTab)
Vue.use(PageLoading)
Vue.use(VueClipboard)
Vue.use(PermissionHelper)
Vue.use(VueCropper)
Vue.use(VueCodeMirror)
Vue.use(contentmenu)