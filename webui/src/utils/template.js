export function extractVars(content) {
    return [...content.matchAll(/{{([\w\._-]+)}}/g)]
        .map((m) => m[1])
        .reduce((acc, val) => Object.assign(acc, { [val]: null }), {})
}