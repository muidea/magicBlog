import { request, config } from 'utils'

const { api } = config
const { summaryQuery, articleCreate, catalogCreate } = api


export async function querySummary(params) {
  return request({
    url: summaryQuery,
    method: 'get',
    data: params,
  })
}

export async function createCatalog(params) {
  return request({
    url: catalogCreate,
    method: 'post',
    data: params,
  })
}

export async function createArticle(params) {
  return request({
    url: articleCreate,
    method: 'post',
    data: params,
  })
}
