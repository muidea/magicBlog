import { request, config } from 'utils'

const { api } = config
const { articlePage } = api

export async function queryArticle(params) {
  return request({
    url: articlePage,
    method: 'get',
    data: params,
  })
}
