import { queryIndex } from 'services/index'
import queryString from 'query-string'

export default {

  namespace: 'index',

  state: {
    summaryList: [],
  },

  subscriptions: {
    setup({ dispatch, history }) {
      history.listen((location) => {
        if (location.pathname === '/') {
          dispatch({
            type: 'queryIndex',
            payload: queryString.parse(location.search),
          })
        }
      })
    },
  },

  effects: {
    *queryIndex({ payload }, { call, put }) {
      const result = yield call(queryIndex, { ...payload })
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
