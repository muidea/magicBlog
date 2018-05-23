import React from 'react'
import { Router, Route, Switch } from 'dva/router'
import { MainLayout } from './components'
import IndexPage from './routes/index'
import CatalogPage from './routes/catalog'
import ContactPage from './routes/contact'
import AboutPage from './routes/about'
import ArticlePage from './routes/article'
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
          <Route component={ErrorPage} />
        </Switch>
      </MainLayout>
    </Router>
  )
}

export default RouterConfig
