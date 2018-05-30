import { routerRedux } from 'dva/router'
import { loginUser, logoutUser } from 'services/maintain'

export default {

  namespace: 'maintain',

  state: {
  },

  subscriptions: {
    setup() {
    },
  },

  effects: {

    *loginUser({ payload }, { call, put }) {
      const result = yield call(loginUser, { ...payload })
      const { data } = result
      console.log(result)

      if (data !== null && data !== undefined) {
        const { errorCode, reason, onlineUser, authToken, sessionID } = data
        if (errorCode === 0) {
          yield put({ type: 'app/save', payload: { isLogin: errorCode === 0, authToken, sessionID, onlineUser } })
          yield put(routerRedux.push({
            pathname: '/',
          }))
        } else {
          throw reason
        }
      }
    },

    *logoutUser({ payload }, { call, put }) {
      const result = yield call(logoutUser, { ...payload })
      const { data } = result
      console.log(result)

      if (data !== null && data !== undefined) {
        const { errorCode, onlineUser } = data

        yield put({ type: 'save', payload: { isLogin: errorCode === 0, onlineUser } })
      }
    },

  },

  reducers: {
    save(state, action) {
      return { ...state, ...action.payload }
    },
  },
}
