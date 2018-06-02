import { request, config } from 'utils'

const { api } = config
const { catalogQuerySummary, catalogQuerySummaryByID } = api

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

