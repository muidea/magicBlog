import { queryStatus, loginUser, logoutUser } from 'services/maintain'
import queryString from 'query-string'

export default {

  namespace: 'maintain',

  state: {
    isLogin: false,
    authToken: '',
    sessionID: '',
    onlineUser: {},
  },

  subscriptions: {
    setup({ dispatch, history }) {
      history.listen((location) => {
        if (location.pathname === '/maintain') {
          dispatch({
            type: 'queryStatus',
            payload: queryString.parse(location.search),
          })
        }
      })
    },
  },

  effects: {
    *queryStatus({ payload }, { call, put, select }) {
      const { authToken, sessionID } = yield select(_ => _.maintain)
      const result = yield call(queryStatus, { authToken, sessionID })
      const { data } = result

      if (data !== null && data !== undefined) {
        const { errorCode, onlineUser } = data

        yield put({ type: 'save', payload: { isLogin: errorCode === 0, onlineUser } })
      }
    },

    *loginUser({ payload }, { call, put }) {
      const result = yield call(loginUser, { ...payload })
      const { data } = result
      console.log(result)

      if (data !== null && data !== undefined) {
        const { errorCode, onlineUser } = data

        yield put({ type: 'save', payload: { isLogin: errorCode === 0, onlineUser } })
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
