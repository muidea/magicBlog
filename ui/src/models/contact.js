import { queryContact } from 'services/contact'
import queryString from 'query-string'

export default {

  namespace: 'contact',

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
        if (location.pathname === '/contact') {
          dispatch({
            type: 'queryContact',
            payload: queryString.parse(location.search),
          })
        }
      })
    },
  },

  effects: {
    *queryContact({ payload }, { call, put }) {
      const result = yield call(queryContact, { ...payload })
      const { data } = result
      console.log(result)
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
