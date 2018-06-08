import { request, config } from 'utils'

const { api } = config
const { summaryQuery, articleCreate, catalogCreate, catalogDelete, articleDelete } = api


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

export async function deleteCatalog(params) {
  return request({
    url: catalogDelete,
    method: 'delete',
    data: params,
  })
}

export async function deleteArticle(params) {
  return request({
    url: articleDelete,
    method: 'delete',
    data: params,
  })
}
