import { user } from '../stores/auth'
import { errors } from '../stores/alerts'
import { sanitize } from '../utils/url'
import { config } from '../stores/config'

const basePath = config.get().basePath

function apiUrl() {
  return sanitize(`/${basePath}/${process?.env?.apiUrl || '/api'}`)
}

function baseHeaders() {
  return {
    'Authorization': `Basic ${ user.getToken() }`,
    'Accept': 'application/json',
    'Content-Type': 'application/json'
  }
}

async function request(path, data, options) {
  options.headers = { ...baseHeaders(), ...(options.headers || {}) }
  if (data) {
    options = { ...options, ...{ body: JSON.stringify(data) } }
  }
  const response = await fetch(`${ apiUrl() }${ path.startsWith('/') ? '' : '/' }${ path }`, options)

  if (response.status === 401) {
    user.logout()
    if (!options.skipRedirect) {
      window.location.replace('')
    }
  }

  if (response.status >= 400) {
    errors.spawn(`Error (${ response.status }): ${ await response.text() }`)
    throw new Error(`request failed with status ${ response.status }`)
  }

  return { data: options.raw ? response.text() : response.json(), response: response }
}

export { apiUrl, baseHeaders, request }
