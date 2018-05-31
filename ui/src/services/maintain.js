import { request, config } from 'utils'

const { api } = config
const { summaryQuery } = api

export async function querySummary(params) {
  return request({
    url: summaryQuery,
    method: 'get',
    data: params,
  })
}
