import axios from 'axios'

const response = function (response) {
  return response
}

const error = function (error) {
  return Promise.reject(error)
}

axios.interceptors.response.use(response, error)

const instance = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json'
  },
  withCredentials: true
})

export default instance
