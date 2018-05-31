import { request, config } from 'utils'

const { api } = config
const { contactQuery } = api

export async function queryContact(params) {
  return request({
    url: contactQuery,
    method: 'get',
    data: params,
  })
}
