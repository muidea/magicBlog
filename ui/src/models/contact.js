import { queryContact } from 'services/contact'
import qs from 'qs'

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
            payload: qs.parse(location.search),
          })
        }
      })
    },
  },

  effects: {
    *queryContact({ payload }, { call, put }) {
      const result = yield call(queryContact, { ...payload })
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
