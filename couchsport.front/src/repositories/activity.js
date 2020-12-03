import repo from './repository'

export default {
  all() {
    return repo.get('/activities')
  }
}
