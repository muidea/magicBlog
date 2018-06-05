import { routerRedux } from 'dva/router'
import queryString from 'query-string'
import { querySummary, createCatalog, createArticle } from '../services/maintain'
import { queryCatalogSummary } from '../services/catalog'
import { queryArticle } from '../services/article'

export default {
  namespace: 'maintain',

  state: {
    itemList: [],
    action: { command: '', id: -1, type: '', name: '' },
  },

  subscriptions: {
    setup({ dispatch, history }) {
      history.listen((location) => {
        if (location.pathname === '/maintain') {
          dispatch({
            type: 'refreshContent',
            payload: queryString.parse(location.search),
          })
        }
      })
    },
  },

  effects: {
    *redirectContent({ payload }, { put }) {
      const { url } = payload
      yield put(routerRedux.push({
        pathname: url,
      }))
    },

    *refreshContent({ payload }, { call, put, select }) {
      const { isLogin, authToken, sessionID } = yield select(_ => _.app)
      if (!isLogin) {
        yield put(routerRedux.push({
          pathname: '/login',
        }))
        return
      }

      console.log(payload)
      {
        const summaryResult = yield call(querySummary, { authToken, sessionID })
        const { data } = summaryResult
        const { errorCode, reason, itemList } = data
        if (errorCode === 0) {
          yield put({ type: 'save', payload: { itemList } })
        } else {
          throw reason
        }
      }

      const { command, id, name, type } = payload
      if (command === 'view') {
        if (type === 'article') {
          const articleResult = yield call(queryArticle, { id })
          const { data } = articleResult
          const { errorCode, reason, content } = data
          if (errorCode === 0) {
            yield put({ type: 'save', payload: { action: { command, id, name, type, data: content } } })
          } else {
            throw reason
          }
        } else {
          const summaryResult = yield call(queryCatalogSummary, { id })
          const { data } = summaryResult
          const { errorCode, reason, summaryList } = data
          if (errorCode === 0) {
            yield put({ type: 'save', payload: { action: { command, id, name, type, data: summaryList } } })
          } else {
            throw reason
          }
        }
      } else if (command === 'add') {
        yield put({ type: 'save', payload: { action: { command, id, name, type, data: {} } } })
      } else if (command === 'modify') {
        yield put({ type: 'save', payload: { action: { command, id, name, type, data: {} } } })
      } else {
      }
    },


    *submitContent({ payload }, { call, put, select }) {
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
