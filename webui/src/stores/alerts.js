import { writable } from 'svelte/store'

function createErrors() {
    const { subscribe, update } = writable([])
    const spawn = (message) => {
        update(alerts => alerts = [...alerts, message])
        setTimeout(() => update(alerts => alerts = alerts.filter(a => a !== message)), 3000)
    }
    return { subscribe, spawn }
}

function createInfos() {
    const { subscribe, update } = writable([])
    const spawn = (message) => {
        update(alerts => alerts = [...alerts, message])
        setTimeout(() => update(alerts => alerts = alerts.filter(a => a !== message)), 3000)
    }
    return { subscribe, spawn }
}

function createSuccesses() {
    const { subscribe, update } = writable([])
    const spawn = (message) => {
        update(alerts => alerts = [...alerts, message])
        setTimeout(() => update(alerts => alerts = alerts.filter(a => a !== message)), 3000)
    }
    return { subscribe, spawn }
}

export const errors = createErrors()
export const infos = createInfos()
export const successes = createSuccesses()