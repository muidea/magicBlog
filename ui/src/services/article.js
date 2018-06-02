import { request, config } from 'utils'

const { api } = config
const { articleQuery } = api

export async function queryArticle(params) {
  return request({
    url: articleQuery,
    method: 'get',
    data: params,
  })
}

