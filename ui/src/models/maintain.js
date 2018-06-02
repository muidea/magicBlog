import { routerRedux } from 'dva/router'
import queryString from 'query-string'
import { querySummary, createCatalog, createArticle } from '../services/maintain'
import { queryCatalogSummary } from '../services/catalog'
import { queryArticle } from '../services/article'

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
      const { isLogin, authToken, sessionID } = yield select(_ => _.app)
      if (!isLogin) {
        yield put(routerRedux.push({
          pathname: '/login',
        }))
        return
      }

      if (authToken) {
        payload = { ...payload, authToken }
      }
      if (sessionID) {
        payload = { ...payload, sessionID }
      }
      const result = yield call(querySummary, { ...payload })
      const { data } = result
      const { errorCode, reason, summaryList } = data
      if (errorCode === 0) {
        yield put({ type: 'save', payload: { summaryList } })
      } else {
        throw reason
      }
    },

    *querySelectContent({ payload }, { call, put }) {
      const { id, type } = payload
      if (type === 'catalog') {
        const result = yield call(queryCatalogSummary, { id })
        const { data } = result
        yield put({ type: 'save', payload: { action: { type: 'viewContent', value: { content: data, currentItem: { ...payload } } } } })
      } else if (type === 'article') {
        const result = yield call(queryArticle, { id })
        const { data } = result
        yield put({ type: 'save', payload: { action: { type: 'viewContent', value: { content: data, currentItem: { ...payload } } } } })
      } else {
        throw type
      }
    },

    *addCatalog({ payload }, { put }) {
      yield put({ type: 'save', payload: { action: { type: 'addCatalog', value: { name: '', description: '', parent: { ...payload } } } } })
    },

    *submitCatalog({ payload }, { call, put, select }) {
      const { authToken, sessionID } = yield select(_ => _.app)

      if (authToken) {
        payload = { ...payload, authToken }
      }
      if (sessionID) {
        payload = { ...payload, sessionID }
      }
      const result = yield call(createCatalog, { ...payload })
      const { data } = result
      const { errorCode, reason } = data
      if (errorCode === 0) {
        yield put({ type: 'save', payload: { action: { type: 'viewContent', value: { data: {}, currentItem: { } } } } })
      } else {
        throw reason
      }
    },

    *addArticle({ payload }, { put }) {
      yield put({ type: 'save', payload: { action: { type: 'addArticle', value: { title: '', content: '', parent: { ...payload } } } } })
    },

    *submitArticle({ payload }, { call, put, select }) {
      const { authToken, sessionID } = yield select(_ => _.app)

      if (authToken) {
        payload = { ...payload, authToken }
      }
      if (sessionID) {
        payload = { ...payload, sessionID }
      }
      const result = yield call(createArticle, { ...payload, authToken, sessionID })
      const { data } = result
      const { errorCode, reason } = data
      if (errorCode === 0) {
        yield put(routerRedux.push({
          pathname: '/maintain',
        }))
      } else {
        throw reason
      }
    },
  },

  reducers: {
    save(state, action) {
      return { ...state, ...action.payload }
    },
  },
}
