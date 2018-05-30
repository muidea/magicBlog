export default {

  namespace: 'maintain',

  state: {
  },

  subscriptions: {
    setup() {
    },
  },

  effects: {

  },

  reducers: {
    save(state, action) {
      return { ...state, ...action.payload }
    },
  },
}
