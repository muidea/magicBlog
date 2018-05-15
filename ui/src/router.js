import React from 'react'
import { Router, Route, Switch } from 'dva/router'
import IndexPage from './routes/IndexPage'
import CatalogPage from './routes/CatalogPage'

function RouterConfig({ history }) {
  return (
    <Router history={history}>
      <Switch>
        <Route path="/" exact component={IndexPage} />
        <Route path="/catalogs" exact component={CatalogPage} />
      </Switch>
    </Router>
  )
}

export default RouterConfig
