import { queryArticle } from 'services/article'
import pathToRegexp from 'path-to-regexp'

export default {

  namespace: 'article',

  state: {
    name: '',
    creater: { id: 0, name: '' },
    createDate: '',
    catalog: [],
    content: '',
  },

  subscriptions: {
    setup({ dispatch, history }) {
      history.listen((location) => {
        const match = pathToRegexp('/article/:i').exec(location.pathname)
        if (match) {
          dispatch({ type: 'queryArticle', payload: { id: match[1] } })
        }
      })
    },
  },

  effects: {
    *queryArticle({ payload }, { call, put }) {
      const result = yield call(queryArticle, { ...payload })
      const { data } = result
      const { content } = data
      yield put({ type: 'save', payload: { ...content } })
    },
  },

  reducers: {
    save(state, action) {
      return { ...state, ...action.payload }
    },
  },
}
