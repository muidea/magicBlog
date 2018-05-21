import React from 'react'
import { connect } from 'dva'
import styles from './index.less'

function ErrorPage() {
  return (
    <div className={styles.normal}>
      <h1 className={styles.title}>ErrorPage</h1>
    </div>
  )
}

ErrorPage.propTypes = {}

export default connect()(ErrorPage)
