import { request, config } from 'utils'

const { api } = config
const { aboutQuery } = api

export async function queryAbout(params) {
  return request({
    url: aboutQuery,
    method: 'get',
    data: params,
  })
}
