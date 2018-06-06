import { request, config } from 'utils'

const { api } = config
const { catalogQuerySummary, catalogQuerySummaryByID, catalogQuery } = api

export async function queryCatalogSummary(params) {
  return request({
    url: catalogQuerySummary,
    method: 'get',
    data: params,
  })
}

export async function queryCatalogSummaryByID(params) {
  return request({
    url: catalogQuerySummaryByID,
    method: 'get',
    data: params,
  })
}

export async function queryCatalogByID(params) {
  return request({
    url: catalogQuery,
    method: 'get',
    data: params,
  })
}
