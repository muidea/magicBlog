import { request, config } from 'utils'

const { api } = config
const { indexPage } = api

export async function queryIndex() {
  return request(indexPage)
}
