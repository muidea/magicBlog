import { request, config } from 'utils'

const { api } = config
const { summaryList } = api

export async function querySummary(params) {
  return request({
    url: summaryList,
    method: 'get',
    data: params,
  })
}
