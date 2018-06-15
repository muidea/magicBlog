import { queryNoFound } from 'services/error'
import qs from 'qs'

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
        if (location.pathname === '/404') {
          dispatch({
            type: 'queryNoFound',
            payload: qs.parse(location.search),
          })
        }
      })
    },
  },

  effects: {
    *queryNoFound({ payload }, { call, put }) {
      const result = yield call(queryNoFound, { ...payload })
      const { data } = result
      yield put({ type: 'save', payload: { ...data } })
    },
  },

  reducers: {
    save(state, action) {
      return { ...state, ...action.payload }
    },
  },
}
