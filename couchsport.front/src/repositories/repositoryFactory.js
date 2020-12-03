import userRepository from './user'
import activityRepository from './activity'

const repositories = {
  user: userRepository,
  activity: activityRepository
}

export const repositoryFactory = {
  get: (name) => repositories[name]
}
