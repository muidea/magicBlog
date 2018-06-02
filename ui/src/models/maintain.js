import { routerRedux } from 'dva/router'
import queryString from 'query-string'
import { querySummary } from 'services/maintain'
import { queryCatalog, createCatalog } from 'services/catalog'
import { queryArticle, createArticle } from 'services/article'

export default {

  namespace: 'maintain',

  state: {
    summaryList: [],
    action: { type: 'viewContent', value: { data: {}, currentItem: { } } },
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
        yield put({ type: 'save', payload: { action: { type: 'viewContent', value: { data, currentItem: { ...payload } } } } })
      } else if (type === 'article') {
        const result = yield call(queryArticle, { id })
        const { data } = result
        yield put({ type: 'save', payload: { action: { type: 'viewContent', value: { data, currentItem: { ...payload } } } } })
      }
    },

    *addCatalog({ payload }, { put }) {
      yield put({ type: 'save', payload: { action: { type: 'addCatalog', value: { name: '', description: '', parent: { ...payload } } } } })
    },

    *submitCatalog({ payload }, { call, select }) {
      const { authToken, sessionID } = yield select(_ => _.app)

      yield call(createCatalog, { ...payload, authToken, sessionID })
    },

    *addArticle({ payload }, { put }) {
      yield put({ type: 'save', payload: { action: { type: 'addArticle', value: { title: '', content: '', parent: { ...payload } } } } })
    },

    *submitArticle({ payload }, { call, select }) {
      const { authToken, sessionID } = yield select(_ => _.app)

      yield call(createArticle, { ...payload, authToken, sessionID })
    },

  },

  reducers: {
    save(state, action) {
      return { ...state, ...action.payload }
    },
  },
}
