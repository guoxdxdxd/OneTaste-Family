import request from './index'

export function fetchDishes(params = {}) {
  return request.get('/dishes', { params })
}

export function createDish(payload) {
  return request.post('/dishes', payload)
}

export function updateDish(id, payload) {
  return request.put(`/dishes/${id}`, payload)
}

export function deleteDish(id) {
  return request.delete(`/dishes/${id}`)
}

export function getDishDetail(id) {
  return request.get(`/dishes/${id}`)
}
