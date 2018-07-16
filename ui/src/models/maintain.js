import { routerRedux } from 'dva/router'
import { querySummary, createCatalog, createArticle, updateArticle, updateCatalog, deleteArticle, deleteCatalog } from '../services/maintain'
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
            const { errorCode, reason, summary } = data
            if (errorCode === 0) {
              yield put({ type: 'save', payload: { action: { ...payload, data: summary } } })
            } else {
              throw reason
            }
          } else {
            const summaryResult = yield call(queryCatalogSummaryByID, { id })
            const { data } = summaryResult
            const { errorCode, reason, summary } = data
            if (errorCode === 0) {
              yield put({ type: 'save', payload: { action: { ...payload, data: summary } } })
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
      } else if (command === 'delete') {
        if (type === 'article') {
          const articleResult = yield call(deleteArticle, { id, authToken, sessionID })
          const { data } = articleResult
          const { errorCode, reason } = data
          if (errorCode === 0) {
            yield put({ type: 'save', payload: { action: { command: 'view' } } })
          } else {
            throw reason
          }
        } else if (type === 'catalog') {
          const catalogResult = yield call(deleteCatalog, { id, authToken, sessionID })
          const { data } = catalogResult
          const { errorCode, reason } = data
          if (errorCode === 0) {
            yield put({ type: 'save', payload: { action: { command: 'view' } } })
          } else {
            throw reason
          }
        }
      }
    },

    *submitContent({ payload }, { call, put, select }) {
      const { type, command } = payload
      const { authToken, sessionID } = yield select(_ => _.app)

      if (authToken) {
        payload = { ...payload, authToken }
      }
      if (sessionID) {
        payload = { ...payload, sessionID }
      }

      delete payload.command

      if (type === 'article') {
        const serviceFunc = (command === 'add') ? createArticle : updateArticle
        const articleResult = yield call(serviceFunc, { ...payload })
        const { data } = articleResult
        const { errorCode, reason, content } = data
        if (errorCode === 0) {
          yield put({ type: 'save', payload: { action: { command: 'view', ...content } } })
        } else {
          throw reason
        }
      } else if (type === 'catalog') {
        const serviceFunc = (command === 'add') ? createCatalog : updateCatalog
        const catalogResult = yield call(serviceFunc, { ...payload })
        const { data } = catalogResult
        const { errorCode, reason, content } = data
        if (errorCode === 0) {
          yield put({ type: 'save', payload: { action: { command: 'view', ...content } } })
        } else {
          throw reason
        }
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
