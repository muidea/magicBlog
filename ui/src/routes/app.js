import React from 'react'
import PropTypes from 'prop-types'
import { connect } from 'dva'
import { FrameLayout } from 'components'

function App({ children, dispatch, app, loading, history }) {
  const { location } = history
  return (
    <FrameLayout location={location}>
      <div>
        {children}
      </div>
    </FrameLayout>    
  )
}

App.propTypes = {
  children: PropTypes.element.isRequired,
  history: PropTypes.object,
  dispatch: PropTypes.func,
  app: PropTypes.object,
  loading: PropTypes.object,
}

export default connect(({ app, loading }) => ({ app, loading }))(App)
