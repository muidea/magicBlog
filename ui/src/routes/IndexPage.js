import React from 'react'
import PropTypes from 'prop-types'
import { connect } from 'dva'
import styles from './IndexPage.css'

function IndexPage({index}) {
  const { msg } = index
  return (
    <div className={styles.normal}>
      <h1 className={styles.title}>IndexPage</h1>
      <p>{msg}</p>
    </div>
  )
}
 
IndexPage.propTypes = {
  index: PropTypes.object,
}

export default connect(({ index }) => ({ index }))(IndexPage)

