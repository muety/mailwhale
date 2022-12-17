function basicAuthEncode(user, password) {
  return btoa(`${ user }:${ password }`)
}

function basicAuthDecode(token) {
  if (!token) return null
  return atob(token).split(':')
}

export { basicAuthEncode, basicAuthDecode }
