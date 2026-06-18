import axios from 'axios'

const API_URL = 'http://localhost:8080/api' // Jika jalan untuk laptop sendiri uncomment url ini.

// const API_URL = 'http://192.168.1.5:8080/api' // IP Laptop atau IP Server agar bisa diakses dari laptop lain di jaringan yang sama

const api = axios.create({ baseURL: API_URL })

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

api.interceptors.response.use(
  (res) => res,
  (err) => {
    if (err.response?.status === 401) {
      localStorage.clear()
      window.location.href = '/login'
    }
    return Promise.reject(err)
  }
)

export default api
