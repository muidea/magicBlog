import { message } from 'antd'
import dva from 'dva'
import createLoading from 'dva-loading'
import createHistory from 'history/createBrowserHistory'
import './index.css'

// 1. Initialize
const app = dva({
  ...createLoading({ effects: true }),
  history: createHistory(),
  onError(error) {
    message.error(error.message)
  },
})

// 2. Plugins
// app.use({});

// 3. Model
app.model(require('models/index'))
app.model(require('models/catalog'))
app.model(require('models/contact'))
app.model(require('models/about'))
app.model(require('models/article'))
app.model(require('models/error'))
app.model(require('models/maintain'))

// 4. Router
app.router(require('./router'))

// 5. Start
app.start('#root')
