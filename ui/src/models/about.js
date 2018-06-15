import { queryAbout } from 'services/about'
import qs from 'qs'

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
            payload: qs.parse(location.search),
          })
        }
      })
    },
  },

  effects: {
    *queryAbout({ payload }, { call, put }) {
      const result = yield call(queryAbout, { ...payload })
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
