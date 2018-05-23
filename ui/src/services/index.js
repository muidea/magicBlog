import { request, config } from 'utils'

const { api } = config
const { indexPage } = api

export async function queryIndex(params) {
  return request({
    url: indexPage,
    method: 'get',
    data: params,
  })
}
