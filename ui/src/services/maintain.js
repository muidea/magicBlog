import { request, config } from 'utils'

const { api } = config
const { maintainStatus, maintainLogin, maintainLogout } = api

export async function queryStatus(params) {
  return request({
    url: maintainStatus,
    method: 'get',
    data: params,
  })
}

export async function loginUser(params) {
  return request({
    url: maintainLogin,
    method: 'post',
    data: params,
  })
}

export async function logoutUser(params) {
  return request({
    url: maintainLogout,
    method: 'delete',
    data: params,
  })
}
