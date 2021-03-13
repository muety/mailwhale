import { writable } from 'svelte/store'
import { basicAuthEncode, basicAuthDecode } from '../utils/auth'

// https://svelte.dev/tutorial/custom-stores

function createUser() {
    const { subscribe, set } = writable(null)

    let token = null

    const load = () => {
        const t = sessionStorage.getItem('auth_token')
        set((basicAuthDecode(t) || [null])[0])
        token = t
    }
    const save = (t) => {
        if (!t) sessionStorage.removeItem('auth_token')
        else sessionStorage.setItem('auth_token', t)
        token = t
    }
    const login = ({ username, password }) => {
        set(username)
        const t = basicAuthEncode(username, password)
        save(t)
    }
    const logout = () => save(null)

    const getToken = () => token

    return { subscribe, load, login, logout, getToken }
}

export const user = createUser()