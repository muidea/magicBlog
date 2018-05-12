import React from 'react'
import { Router, Route, Switch } from 'dva/router'
import IndexPage from './routes/IndexPage'

import Catalogs from './routes/Catalogs.js'

function RouterConfig({ history }) {
  return (
    <Router history={history}>
      <Switch>
        <Route path="/catalogs" component={Catalogs} />
        <Route path="/" component={IndexPage} />
      </Switch>
    </Router>
  )
}

export default RouterConfig
