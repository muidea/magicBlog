import fetch from 'dva/fetch'

function parseJSON(response) {
  return response.json()
}

function checkStatus(response) {
  if (response.status >= 200 && response.status < 300) {
    return response
  }

  const error = new Error(response.statusText)
  error.response = response
  throw error
}

/**
 * Requests a URL, returning a promise.
 *
 * @param  {string} url       The URL we want to request
 * @param  {object} [options] The options we want to pass to "fetch"
 * @return {object}           An object containing either "data" or "err"
 */
export default function request(options) {
  if (options.url) {
    if (options.data) {
      const { id } = options.data
      let { url } = options
      if (id !== undefined) {
        delete options.data.id
        url = url.replace(':id', id)
      }

      options = {
        ...options,
        url,
      }
    }
  }

  const { url, method } = options
  if (method === 'post') {
    options = {
      body: JSON.stringify(options.data),
      cache: 'no-cache',
      headers: {
        'user-agent': 'Mozilla/4.0 MDN Example',
        'content-type': 'application/json',
      },
      method: 'POST',
      mode: 'cors',
    }
  }

  return fetch(url, options)
    .then(checkStatus)
    .then(parseJSON)
    .then(data => ({ data }))
    .catch(err => ({ err }))
}
