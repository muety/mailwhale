import { user } from '../stores/auth'

function apiUrl() {
    return process.env.apiUrl || '/api'
}

function baseHeaders() {
    return {
        'Authorization': `Basic ${user.getToken()}`,
        'Accept': 'application/json',
        'Content-Type': 'application/json'
    }
}

async function request(path, data, options) {
    options.headers = { ...baseHeaders(), ...(options.headers || {}) }
    if (data) {
        options = { ...options, ...{ body: JSON.stringify(data) } }
    }
    const response = await fetch(`${apiUrl()}${path.startsWith('/') ? '' : '/'}${path}`, options)

    if (response.status === 401) {
        user.logout()
        window.location.replace('/')
    }

    if (response.status >= 400) {
        throw new Error(`request faield with status ${response.status}`)
    }

    return { data: response.json(), response: response }
}

export { apiUrl, baseHeaders, request }