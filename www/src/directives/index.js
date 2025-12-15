/**
 * 全局指令注册
 * Vue3 版本
 */
import action from './action'

/**
 * 注册全局指令
 * @param {App} app Vue3 应用实例
 */
export function setupDirectives(app) {
  // 注册 v-action 指令
  app.directive('action', action)
}
