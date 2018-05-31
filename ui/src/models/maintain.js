import { routerRedux } from 'dva/router'
import queryString from 'query-string'
import { querySummary } from 'services/maintain'
import { queryCatalog } from 'services/catalog'
import { queryArticle } from 'services/article'

export default {

  namespace: 'maintain',

  state: {
    summaryList: [],
    currentSelect: { summary: {}, content: {} },
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

    *querySelectContent({ payload }, { call, put }) {
      const { id, type } = payload
      if (type === 'catalog') {
        const result = yield call(queryCatalog, { id })
        const { data } = result
        yield put({ type: 'save', payload: { currentSelect: { summary: payload, content: data } } })
      } else if (type === 'article') {
        const result = yield call(queryArticle, { id })
        const { data } = result
        yield put({ type: 'save', payload: { currentSelect: { summary: payload, content: data } } })
      }
    },
  },

  reducers: {
    save(state, action) {
      return { ...state, ...action.payload }
    },
  },
}
