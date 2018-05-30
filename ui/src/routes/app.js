/* global window */
import React from 'react'
import PropTypes from 'prop-types'
import { connect } from 'dva'
import { Helmet } from 'react-helmet'
import { withRouter } from 'dva/router'
import { MainLayout, MaintainLayout } from '../components'

const App = ({ children, app, history, dispatch, loading }) => {
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
      <MainLayout history={history} dispatch={dispatch} loading={loading}>
        { children }
      </MainLayout>
      }
      { isLogin &&
      <MaintainLayout history={history} user={onlineUser} logoutHandler={onLogoutHandler} dispatch={dispatch} loading={loading}>
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
  loading: PropTypes.object,
}

export default withRouter(connect(({ app, loading }) => ({ app, loading }))(App))
