import React from 'react'
import { Router, Route, Switch } from 'dva/router'
import { MainLayout } from './components'
import IndexPage from './routes/index'
import CatalogPage from './routes/catalog'
import ContactPage from './routes/contact'
import AboutPage from './routes/about'
import ContentPage from './routes/content'
import ErrorPage from './routes/error'

function RouterConfig({ history }) {
  return (
    <Router history={history}>
      <MainLayout history={history}>
        <Switch>
          <Route exact path="/" component={IndexPage} />
          <Route exact path="/catalog" component={CatalogPage} />
          <Route exact path="/contact" component={ContactPage} />
          <Route exact path="/about" component={AboutPage} />
          <Route exact path="/content" component={ContentPage} />
          <Route component={ErrorPage} />
        </Switch>
      </MainLayout>
    </Router>
  )
}

export default RouterConfig
