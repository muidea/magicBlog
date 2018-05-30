import { request, config } from 'utils'

const { api } = config
const { userStatus, userLogin, userLogout } = api

export async function queryStatus(params) {
  return request({
    url: userStatus,
    method: 'get',
    data: params,
  })
}

export async function loginUser(params) {
  return request({
    url: userLogin,
    method: 'post',
    data: params,
  })
}

export async function logoutUser(params) {
  return request({
    url: userLogout,
    method: 'delete',
    data: params,
  })
}
