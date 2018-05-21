import React from 'react'
import { connect } from 'dva'
import styles from './index.less'

function ContentPage() {
  return (
    <div className={styles.normal}>
      <h1 className={styles.title}>ContentPage</h1>
    </div>
  )
}

ContentPage.propTypes = {}

export default connect()(ContentPage)
