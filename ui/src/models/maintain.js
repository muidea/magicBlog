import { queryMaintain } from 'services/maintain'
import queryString from 'query-string'

export default {

  namespace: 'maintain',

  state: {
    isLogin: false,
  },

  subscriptions: {
    setup({ dispatch, history }) {
      history.listen((location) => {
        if (location.pathname === '/maintainA') {
          dispatch({
            type: 'queryMaintain',
            payload: queryString.parse(location.search),
          })
        }
      })
    },
  },

  effects: {
    *queryMaintain({ payload }, { call, put }) {
      const result = yield call(queryMaintain, { payload })
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
