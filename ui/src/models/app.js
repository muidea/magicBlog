/* global window */

import queryString from 'query-string'
import { queryStatus } from 'services/maintain'
import { config } from 'utils'

const { prefix } = config

export default {
  namespace: 'app',
  state: {
    isLogin: true,
    sessionID: window.localStorage.getItem(`${prefix}SessionID`),
    authToken: window.localStorage.getItem(`${prefix}AuthToken`),
    onlineUser: {},
  },

  subscriptions: {
    setupHistory({ dispatch, history }) {
      history.listen((location) => {
        dispatch({
          type: 'queryStatus',
          payload: {
            locationPathname: location.pathname,
            locationQuery: queryString.parse(location.search),
          },
        })
      })
    },
  },

  effects: {
    *queryStatus({ payload }, { call, put, select }) {
      const { authToken, sessionID } = yield select(_ => _.app)
      const result = yield call(queryStatus, { authToken, sessionID })
      const { data } = result

      if (data !== null && data !== undefined) {
        const { errorCode, onlineUser } = data

        yield put({ type: 'save', payload: { isLogin: errorCode === 0, onlineUser } })
      }
    },
  },

  reducers: {
    save(state, { payload }) {
      const { sessionID, authToken } = payload

      if (sessionID && authToken) {
        window.localStorage.setItem(`${prefix}SessionID`, sessionID)
        window.localStorage.setItem(`${prefix}AuthToken`, authToken)
      }

      return { ...state, ...payload }
    },
  },
}
