import { request, config } from 'utils'

const { api } = config
const { noFoundPage } = api

export async function queryNoFound(params) {
  return request({
    url: noFoundPage,
    method: 'get',
    data: params,
  })
}
