import { request, config } from 'utils'

const { api } = config
const { indexQuery } = api

export async function queryIndex(params) {
  return request({
    url: indexQuery,
    method: 'get',
    data: params,
  })
}
