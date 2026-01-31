import request from './index'

/**
 * 创建菜单
 * @param {Object} payload - { date: 'YYYY-MM-DD', meal_type: 'breakfast|lunch|dinner', dish_ids: [] }
 */
export function createMenu(payload) {
  return request.post('/menus', payload)
}

/**
 * 获取每日菜单
 * @param {string} date - 日期，格式：YYYY-MM-DD
 */
export function getDailyMenu(date) {
  return request.get('/menus/daily', { params: { date } })
}

/**
 * 获取每周菜单
 * @param {string} startDate - 开始日期，格式：YYYY-MM-DD
 */
export function getWeeklyMenu(startDate) {
  return request.get('/menus/weekly', { params: { start_date: startDate } })
}

/**
 * 更新菜单
 * @param {string} id - 菜单ID
 * @param {Object} payload - { date?: 'YYYY-MM-DD', meal_type?: 'breakfast|lunch|dinner', dish_ids?: [] }
 */
export function updateMenu(id, payload) {
  return request.put(`/menus/${id}`, payload)
}

/**
 * 删除菜单
 * @param {string} id - 菜单ID
 */
export function deleteMenu(id) {
  return request.delete(`/menus/${id}`)
}

