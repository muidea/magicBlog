import React from 'react'
import { Router, Route, Redirect, Switch } from 'dva/router'
import { MainLayout } from './components'
import IndexPage from './routes/index'
import CatalogPage from './routes/catalog'
import ContactPage from './routes/contact'
import AboutPage from './routes/about'
import ArticlePage from './routes/article'
import MaintainPage from './routes/maintain'
import ErrorPage from './routes/error'

function RouterConfig({ history }) {
  return (
    <Router history={history}>
      <MainLayout history={history}>
        <Switch>
          <Route exact path="/" component={IndexPage} />
          <Route exact path="/catalog" component={CatalogPage} />
          <Route exact path="/catalog/:id" component={CatalogPage} />
          <Route exact path="/contact" component={ContactPage} />
          <Route exact path="/about" component={AboutPage} />
          <Route exact path="/article/:id" component={ArticlePage} />
          <Route exact path="/maintain" component={MaintainPage} />
          <Route exact path="/404" component={ErrorPage} />
          <Redirect to="/404" />
        </Switch>
      </MainLayout>
    </Router>
  )
}

export default RouterConfig
