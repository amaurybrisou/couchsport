import repo from './repository'

export default {
  get() {
    return repo.get('/profiles/mine')
  },
  update(payload) {
    return repo.post('/profiles/update', payload)
  }
}
