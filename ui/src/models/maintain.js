import { routerRedux } from 'dva/router'
import queryString from 'query-string'
import { querySummary } from 'services/maintain'

export default {

  namespace: 'maintain',

  state: {
    summaryList: [],
  },

  subscriptions: {
    setup({ dispatch, history }) {
      history.listen((location) => {
        if (location.pathname === '/maintain') {
          dispatch({
            type: 'querySummary',
            payload: queryString.parse(location.search),
          })
        }
      })
    },
  },

  effects: {
    *querySummary({ payload }, { call, put, select }) {
      const { isLogin } = yield select(_ => _.app)
      if (!isLogin) {
        yield put(routerRedux.push({
          pathname: '/login',
        }))
        return
      }

      const result = yield call(querySummary, { ...payload })
      const { data } = result
      if (data !== null && data !== undefined) {
        yield put({ type: 'save', payload: { summaryList: data } })
      }
    },
  },

  reducers: {
    save(state, action) {
      return { ...state, ...action.payload }
    },
  },
}
