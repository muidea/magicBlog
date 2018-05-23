import { queryCatalog } from 'services/catalog'
import queryString from 'query-string'

export default {

  namespace: 'catalog',

  state: {
    summaryList: [],
  },

  subscriptions: {
    setup({ dispatch, history }) {
      history.listen((location) => {
        if (location.pathname === '/catalog') {
          dispatch({
            type: 'queryCatalog',
            payload: queryString.parse(location.search),
          })
        }
      })
    },
  },

  effects: {
    *queryCatalog({ payload }, { call, put }) {
      const result = yield call(queryCatalog, { payload })
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
