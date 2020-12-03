import repo from './repository'

export default {
  create(payload) {
    return repo.put('/signup', payload)
  },
  login(payload) {
    return repo.post('/login', payload)
  },
  logout() {
    return repo.get('/logout')
  },
  changePassword(payload) {
    return repo.post('/users/change-password', payload)
  }
}
