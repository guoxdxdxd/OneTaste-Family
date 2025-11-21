import request from './index'

export function searchIngredients(params) {
  return request.get('/ingredients/search', { params })
}

export function fetchIngredientsByCategory(params) {
  return request.get('/ingredients/by-category', { params })
}
