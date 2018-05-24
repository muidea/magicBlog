import { queryNoFound } from 'services/error'
import queryString from 'query-string'

export default {

  namespace: 'error',

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
        dispatch({
          type: 'queryNoFound',
          payload: queryString.parse(location.search),
        })
      })
    },
  },

  effects: {
    *queryNoFound({ payload }, { call, put }) {
      const result = yield call(queryNoFound, { payload })
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
