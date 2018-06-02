/* global window */
import axios from 'axios'
import lodash from 'lodash'
import pathToRegexp from 'path-to-regexp'
import { message } from 'antd'

const fetch = (options) => {
  const {
    method = 'get',
    data,
  } = options
  let {
    url,
  } = options

  const cloneData = lodash.cloneDeep(data)

  try {
    let domin = ''
    if (url.match(/[a-zA-z]+:\/\/[^/]*/)) {
      [domin] = url.match(/[a-zA-z]+:\/\/[^/]*/)
      url = url.slice(domin.length)
    }
    const match = pathToRegexp.parse(url)
    url = pathToRegexp.compile(url)(data)
    for (const item of match) {
      if (item instanceof Object && item.name in cloneData) {
        delete cloneData[item.name]
      }
    }
    url = domin + url
  } catch (e) {
    message.error(e.message)
  }

  switch (method.toLowerCase()) {
    case 'get':
      return axios.get(url, { params: cloneData })
    case 'delete':
      return axios.delete(url, { data: cloneData })
    case 'post':
      return axios.post(url, cloneData)
    case 'put':
      return axios.put(url, cloneData)
    case 'patch':
      return axios.patch(url, cloneData)
    default:
      return axios(options)
  }
}

export default function request(options) {
  if (options.url) {
    if (options.data) {
      const { id, authToken, sessionID } = options.data
      let { url } = options
      if (id !== undefined && (options.method !== 'post')) {
        delete options.data.id
        url = url.replace(':id', id)
      }

      if (authToken !== undefined) {
        delete options.data.authToken
        url = url.concat('?authToken='.concat(authToken))
      }

      if (sessionID !== undefined) {
        delete options.data.sessionID
        url = url.concat('?sessionID='.concat(sessionID))
      }

      options = {
        ...options,
        url,
      }
    }
  }

  console.log(options)

  return fetch(options).then((response) => {
    const { statusText, status } = response
    const { data } = response

    console.log(data)

    return Promise.resolve({
      success: true,
      message: statusText,
      statusCode: status,
      data,
    })
  }).catch((error) => {
    const { response } = error
    let msg
    let statusCode
    if (response && response instanceof Object) {
      const { data, statusText } = response
      statusCode = response.status
      msg = data.message || statusText
    } else {
      statusCode = 600
      msg = error.message || 'Network Error'
    }

    const val = { success: false, statusCode, message: msg }
    return Promise.reject(val)
  })
}
