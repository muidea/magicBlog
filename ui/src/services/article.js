import { request, config } from 'utils'

const { api } = config
const { articleQuery, articleCreate } = api

export async function queryArticle(params) {
  return request({
    url: articleQuery,
    method: 'get',
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
