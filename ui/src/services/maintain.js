import { request, config } from 'utils'

const { api } = config
const { maintainPage } = api

export async function queryIndex(params) {
  return request({
    url: maintainPage,
    method: 'get',
    data: params,
  })
}
