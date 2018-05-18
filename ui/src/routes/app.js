import React from 'react'
import PropTypes from 'prop-types'
import { connect } from 'dva'
import { FrameLayout } from 'components'

function App({ children, location, app }) {
  const { msg } = app
  return (
    <FrameLayout location={location}>
      <div>
        <div>{msg}</div>
        <div>{children}</div>
      </div>
    </FrameLayout>
  )
}

App.propTypes = {
  app: PropTypes.object,
}

export default connect(({ app }) => ({ app }))(App)
