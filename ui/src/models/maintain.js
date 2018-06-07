import { routerRedux } from 'dva/router'
import { querySummary, createCatalog, createArticle } from '../services/maintain'
import { queryCatalogSummary, queryCatalogSummaryByID, queryCatalogByID } from '../services/catalog'
import { queryArticle } from '../services/article'

export default {
  namespace: 'maintain',

  state: {
    itemList: [],
    action: { command: 'view', id: -1, type: 'catalog', name: '' },
  },

  subscriptions: {
    setup({ dispatch, history }) {
      history.listen((location) => {
        if (location.pathname === '/maintain') {
          dispatch({
            type: 'refreshContent',
            payload: { command: 'view', id: -1, type: 'catalog', name: '' },
          })
        }
      })
    },
  },

  effects: {
    *refreshContent({ payload }, { call, put, select }) {
      const { isLogin, authToken, sessionID } = yield select(_ => _.app)
      if (!isLogin) {
        yield put({ type: 'clear' })
        yield put(routerRedux.push({
          pathname: '/login',
        }))
        return
      }

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

      const { command, id, type } = payload
      if (command === 'view') {
        if (type === 'article') {
          const articleResult = yield call(queryArticle, { id })
          const { data } = articleResult
          const { errorCode, reason, content } = data
          if (errorCode === 0) {
            yield put({ type: 'save', payload: { action: { ...payload, data: content } } })
          } else {
            throw reason
          }
        } else if (type === 'catalog') {
          if (id === -1) {
            const summaryResult = yield call(queryCatalogSummary, {})
            const { data } = summaryResult
            const { errorCode, reason, summaryList } = data
            if (errorCode === 0) {
              yield put({ type: 'save', payload: { action: { ...payload, data: summaryList } } })
            } else {
              throw reason
            }
          } else {
            const summaryResult = yield call(queryCatalogSummaryByID, { id })
            const { data } = summaryResult
            const { errorCode, reason, summaryList } = data
            if (errorCode === 0) {
              yield put({ type: 'save', payload: { action: { ...payload, data: summaryList } } })
            } else {
              throw reason
            }
          }
        }
      } else if (command === 'add') {
        yield put({ type: 'save', payload: { action: { ...payload } } })
      } else if (command === 'modify') {
        if (type === 'article') {
          const articleResult = yield call(queryArticle, { id })
          const { data } = articleResult
          const { errorCode, reason, content } = data
          if (errorCode === 0) {
            yield put({ type: 'save', payload: { action: { ...payload, data: content } } })
          } else {
            throw reason
          }
        } else if (type === 'catalog') {
          const catalogResult = yield call(queryCatalogByID, { id, authToken, sessionID })
          const { data } = catalogResult
          const { errorCode, reason, content } = data
          if (errorCode === 0) {
            yield put({ type: 'save', payload: { action: { ...payload, data: content } } })
          } else {
            throw reason
          }
        }
      }
    },

    *submitContent({ payload }, { call, put, select }) {
      const { type } = payload
      const { authToken, sessionID } = yield select(_ => _.app)

      if (authToken) {
        payload = { ...payload, authToken }
      }
      if (sessionID) {
        payload = { ...payload, sessionID }
      }

      if (type === 'article') {
        const articleResult = yield call(createArticle, { ...payload })
        const { data } = articleResult
        const { errorCode, reason, content } = data
        if (errorCode === 0) {
          yield put({ type: 'save', payload: { action: { command: 'view', ...content } } })
        } else {
          throw reason
        }
      } else if (type === 'catalog') {
        const catalogResult = yield call(createCatalog, { ...payload })
        const { data } = catalogResult
        const { errorCode, reason, content } = data
        if (errorCode === 0) {
          yield put({ type: 'save', payload: { action: { command: 'view', ...content } } })
        } else {
          throw reason
        }
      }
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

    clear(state) {
      return { ...state, itemList: [], action: { command: '', id: -1, type: '', name: '' } }
    },
  },
}
