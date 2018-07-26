/* global window */
import React from 'react'
import PropTypes from 'prop-types'
import { connect } from 'dva'
import { Helmet } from 'react-helmet'
import { withRouter } from 'dva/router'
import { config } from 'utils'
import { MainLayout, MaintainLayout } from './common'

const { name } = config

const App = ({ children, app, history, dispatch }) => {
  const { isLogin, onlineUser, authToken, sessionID } = app
  const onLogoutHandler = () => {
    dispatch({ type: 'app/logoutUser', payload: { authToken, sessionID } })
  }

  return (
    <div>
      <Helmet>
        <title>MagicBlog</title>
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
      </Helmet>
      { !isLogin &&
      <MainLayout history={history} declare={name}>
        { children }
      </MainLayout>
      }
      { isLogin &&
      <MaintainLayout history={history} declare={name} user={onlineUser} logoutHandler={onLogoutHandler}>
        { children }
      </MaintainLayout>
      }
    </div>
  )
}

App.propTypes = {
  children: PropTypes.element.isRequired,
  history: PropTypes.object,
  dispatch: PropTypes.func,
  app: PropTypes.object,
}

export default withRouter(connect(({ app, loading }) => ({ app, loading }))(App))
