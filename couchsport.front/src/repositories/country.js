import repo from './repository'

export default {
  all() {
    return repo.get('/languages')
  }
}

// https://restcountries.eu/rest/v2/all
