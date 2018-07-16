import { queryCatalogSummary, queryCatalogSummaryByID } from 'services/catalog'
import qs from 'qs'
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
            payload: qs.parse(location.search),
          })
        } else {
          const match = pathToRegexp('/catalog/:i').exec(location.pathname)
          if (match) {
            if (match) {
              dispatch({ type: 'queryCatalogByID', payload: { id: match[1] } })
            }
          }
        }
      })
    },
  },

  effects: {
    *queryCatalog({ payload }, { call, put }) {
      const result = yield call(queryCatalogSummary, { ...payload })
      const { data } = result
      const { summary } = data
      yield put({ type: 'save', payload: { summaryList: summary } })
    },

    *queryCatalogByID({ payload }, { call, put }) {
      const result = yield call(queryCatalogSummaryByID, { ...payload })
      const { data } = result
      const { summary } = data
      yield put({ type: 'save', payload: { summaryList: summary } })
    },
  },

  reducers: {
    save(state, action) {
      return { ...state, ...action.payload }
    },
  },
}
