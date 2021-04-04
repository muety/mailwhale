import { request } from './api'

async function createUser(signup) {
    return (await request('/user', signup, { method: 'POST' })).data
}

async function updateUser(id, signup) {
    return (await request(`/user/${id}`, signup, { method: 'PUT' })).data
}

export { createUser, updateUser }