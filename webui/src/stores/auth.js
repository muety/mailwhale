import { writable } from 'svelte/store'
import { basicAuthEncode, basicAuthDecode } from '../utils/auth'
import {getMe} from '../api/users'

// https://svelte.dev/tutorial/custom-stores

function createUser() {
    const { subscribe, set } = writable(null)  // store stores the username

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
    const login = async ({ username, password }) => {
        set(username)
        const t = basicAuthEncode(username, password)
        save(t)

        try {
            await getMe({skipRedirect: true})
        } catch (e) {
            logout()
            throw e
        }
    }
    const logout = () => {
        set(null)
        save(null)
    }

    const getToken = () => token

    return { subscribe, load, login, logout, getToken }
}

export const user = createUser()
