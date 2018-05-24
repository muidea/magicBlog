import { queryCatalog, querySingleCatalog } from 'services/catalog'
import queryString from 'query-string'
import pathToRegexp from 'path-to-regexp'

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
        } else {
          const match = pathToRegexp('/catalog/:i').exec(location.pathname)
          if (match) {
            if (match) {
              dispatch({ type: 'querySingleCatalog', payload: { id: match[1] } })
            }
          }
        }
      })
    },
  },

  effects: {
    *queryCatalog({ payload }, { call, put }) {
      const result = yield call(queryCatalog, { ...payload })
      const { data } = result
      if (data !== null && data !== undefined) {
        yield put({ type: 'save', payload: { summaryList: data } })
      }
    },

    *querySingleCatalog({ payload }, { call, put }) {
      const result = yield call(querySingleCatalog, { ...payload })
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
