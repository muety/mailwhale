import { request } from './api'

async function sendMail(mail, asClient) {
    const headers = {}
    if (asClient) {
        headers['X-Client-Id'] = asClient
    }
    return await request('/mail', mail, { method: 'POST', headers })
}

export { sendMail }