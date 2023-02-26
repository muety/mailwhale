const reservedRoutes = new Set([
  'clients',
  'login',
  'signup',
  'clients',
  'mails',
  'imprint',
  'addresses',
  'templates',
])

function createConfig() {
  let parts = location.pathname.split('/')
  let basePath = parts.length >= 2 ? parts[1] : ''
  // if app is served under root path ('/'), but initial entrypoint is a client-side route, e.g. '/clients' then that route would be considered the base path
  // this hacky solution prevents this
  // TODO: use router's API to check the app's client-side root route instead
  if (reservedRoutes.has(basePath)) {
    basePath = ''
  }

  const get = () => ({ basePath })

  return { get }
}

export const config = createConfig()
