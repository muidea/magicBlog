import { request, config } from 'utils'

const { api } = config
const { catalogPage, singleCatalogPage } = api

export async function queryCatalog(params) {
  return request({
    url: catalogPage,
    method: 'get',
    data: params,
  })
}

export async function querySingleCatalog(params) {
  return request({
    url: singleCatalogPage,
    method: 'get',
    data: params,
  })
}
