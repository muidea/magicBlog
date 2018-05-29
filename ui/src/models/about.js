import { queryAbout } from 'services/about'
import queryString from 'query-string'

export default {

  namespace: 'about',

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
        if (location.pathname === '/about') {
          dispatch({
            type: 'queryAbout',
            payload: queryString.parse(location.search),
          })
        }
      })
    },
  },

  effects: {
    *queryAbout({ payload }, { call, put }) {
      const result = yield call(queryAbout, { ...payload })
      const { data } = result
      if (data !== null && data !== undefined) {
        yield put({ type: 'save', payload: { ...data } })
      }
    },
  },

  reducers: {
    save(state, action) {
      return { ...state, ...action.payload }
    },
  },
}
