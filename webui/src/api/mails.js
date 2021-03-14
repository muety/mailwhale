import { request } from './api'

async function sendMail(mail) {
    return await request('/mail', mail, { method: 'POST' })
}

export { sendMail }