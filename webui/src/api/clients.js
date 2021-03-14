import { request } from './api'

async function getClients() {
    return (await request('/client', null, {})).data
}

async function createClient(client) {
    return (await request('/client', client, { method: 'POST' })).data
}

async function deleteClient(id) {
    return await request(`/client/${id}`, null, { method: 'DELETE' })
}

export { getClients, createClient, deleteClient }