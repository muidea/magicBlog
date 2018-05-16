import React from 'react'
import PropTypes from 'prop-types'
import { connect } from 'dva'
import styles from './IndexPage.css'

function IndexPage({location, index}) {
  const { msg } = index
  return (
    <div className={styles.normal}>
      <h1 className={styles.title}>IndexPage</h1>
      <p>{msg}</p>
    </div>
  )
}
 
IndexPage.propTypes = {
  location: PropTypes.object,
  index: PropTypes.object,
}

export default connect(({ index, location }) => ({ index, location }))(IndexPage)

