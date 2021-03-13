import { request } from './api'

async function getClients() {
    return (await request('/client', {})).data
}

export { getClients }