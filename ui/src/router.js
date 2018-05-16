import React from 'react'
import PropTypes from 'prop-types'
import { Router, Route, Switch } from 'dva/router'
import dynamic from 'dva/dynamic'
import App from 'routes/app'

function RouterConfig({ history, app }) {
  const routes = [
    {
      path: '/',
      models: () => [import('models/index')],
      component: () => import('routes/IndexPage'),
    }, {
      path: '/catalog',
      models: () => [import('models/catalog')],
      component: () => import('routes/CatalogPage'),
    }, {
      path: '/content',
      models: () => [import('models/content')],
      component: () => import('routes/ContentPage'),
    },
  ]

  return (
    <Router history={history}>
      <App history={history} app={app}>
        <Switch>
          {
            routes.map(({ path, ...dynamics }, key) => (
              <Route key={key}
                exact
                path={path}
                component={dynamic({
                  app,
                  ...dynamics,
                })}
              />
            ))
          }
        </Switch>
      </App>
    </Router>
  )
}


RouterConfig.propTypes = {
  history: PropTypes.object,
  app: PropTypes.object,
}

export default RouterConfig
