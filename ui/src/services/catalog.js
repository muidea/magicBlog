import { request, config } from 'utils'

const { api } = config
const { catalogQueryAll, catalogQuery } = api

export async function queryCatalog(params) {
  return request({
    url: catalogQueryAll,
    method: 'get',
    data: params,
  })
}

export async function querySingleCatalog(params) {
  return request({
    url: catalogQuery,
    method: 'get',
    data: params,
  })
}
