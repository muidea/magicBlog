import { request, config } from 'utils'

const { api } = config
const { aboutPage } = api

export async function queryAbout(params) {
  return request({
    url: aboutPage,
    method: 'get',
    data: params,
  })
}
