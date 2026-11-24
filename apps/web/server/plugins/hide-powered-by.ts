export default defineNitroPlugin((nitroApp) => {
  nitroApp.hooks.hook('render:response', (response) => {
    delete response.headers['x-powered-by']
    delete response.headers['X-Powered-By']
  })
  nitroApp.hooks.hook('beforeResponse', (event) => {
    removeResponseHeader(event, 'x-powered-by')
    removeResponseHeader(event, 'X-Powered-By')
  })
})
