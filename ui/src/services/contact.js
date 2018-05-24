import { request, config } from 'utils'

const { api } = config
const { contactPage } = api

export async function queryContact(params) {
  return request({
    url: contactPage,
    method: 'get',
    data: params,
  })
}