import React from 'react'
import PropTypes from 'prop-types'
import { Router, Route, Switch } from 'dva/router'
import dynamic from 'dva/dynamic'
import App from 'routes/app'

function RouterConfig({ history, app }) {
  const error = dynamic({
    app,
    component: () => import('./routes/error'),
  })
  const routes = [
    {
      path: '/',
      models: () => [import('models/index')],
      component: () => import('routes/index'),
    }, {
      path: '/catalog',
      models: () => [import('models/catalog')],
      component: () => import('routes/catalog'),
    }, {
      path: '/content',
      models: () => [import('models/content')],
      component: () => import('routes/content'),
    },
  ]

  const { location } = history
  return (
    <Router history={history}>
      <App location={location} app={app}>
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
    </Router>
  )
}


RouterConfig.propTypes = {
  history: PropTypes.object,
  app: PropTypes.object,
}

export default RouterConfig
