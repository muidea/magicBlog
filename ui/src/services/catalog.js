import { request, config } from 'utils'

const { api } = config
const { catalogPage } = api

export async function queryCatalog(params) {
  return request({
    url: catalogPage,
    method: 'get',
    data: params,
  })
}
