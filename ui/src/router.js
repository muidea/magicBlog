import React from 'react'
import PropTypes from 'prop-types'
import { Switch, Route, routerRedux } from 'dva/router'
import dynamic from 'dva/dynamic'
import App from 'routes/app'

const { ConnectedRouter } = routerRedux

const RouterConfig = ({ history, app }) => {
  const error = dynamic({
    app,
    models: () => [import('./models/error')],
    component: () => import('./routes/error/'),
  })
  const routes = [
    {
      path: '/',
      models: () => [import('./models/index')],
      component: () => import('./routes/index/'),
    },
    {
      path: '/catalog',
      models: () => [import('./models/catalog')],
      component: () => import('./routes/catalog/'),
    },
    {
      path: '/catalog/:id',
      models: () => [import('./models/catalog')],
      component: () => import('./routes/catalog/'),
    },
    {
      path: '/contact',
      models: () => [import('./models/contact')],
      component: () => import('./routes/contact/'),
    },
    {
      path: '/about',
      models: () => [import('./models/about')],
      component: () => import('./routes/about/'),
    },
    {
      path: '/article/:id',
      models: () => [import('./models/article')],
      component: () => import('./routes/article/'),
    },
    {
      path: '/maintain',
      models: () => [import('./models/maintain')],
      component: () => import('./routes/maintain/'),
    },
  ]

  return (
    <ConnectedRouter history={history}>
      <App>
        <Switch>
          {
            routes.map(({ path, ...dynamics }, key) => (
              <Route
                key={key}
                exact
                path={path}
                component={dynamic({
                  app,
                  ...dynamics,
                })}
              />
            ))
          }
          <Route component={error} />
        </Switch>
      </App>
    </ConnectedRouter>
  )
}

RouterConfig.propTypes = {
  history: PropTypes.object,
  app: PropTypes.object,
}

export default RouterConfig
