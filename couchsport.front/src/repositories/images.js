import repo from './repository'

export default {
  delete(payload) {
    return repo.post('/images/delete', payload)
  }
}
