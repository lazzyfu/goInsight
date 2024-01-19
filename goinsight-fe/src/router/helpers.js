import _ from 'lodash'

import router from './index.js'

const redirect = (location, replace = false) => {
  const route = router.currentRoute
  const target = {
    name: route.name,
    params: _.cloneDeep(route.params),
    query: _.cloneDeep(route.query),
  }
  _.merge(target, location)
  if (replace) {
    router.replace(target)
  } else {
    router.push(target)
  }
}

export { redirect }
