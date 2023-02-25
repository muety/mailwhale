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
  if (reservedRoutes.has(basePath)) {
    basePath = ''
  }

  const get = () => ({ basePath })

  return { get }
}

export const config = createConfig()
