import { request } from './api'

async function createUser(signup) {
    return (await request('/user', signup, { method: 'POST' })).data
}

async function getMe() {
    return (await request('/user/me', null, {})).data
}

async function updateMe(data) {
    return (await request(`/user/me`, data, { method: 'PUT' })).data
}

export { createUser, updateMe, getMe }