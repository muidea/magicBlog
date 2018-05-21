import React from 'react'
import { connect } from 'dva'
import styles from './index.less'

function AboutPage() {
  return (
    <div className={styles.normal}>
      <h1 className={styles.title}>AboutPage</h1>
    </div>
  )
}

AboutPage.propTypes = {}

export default connect()(AboutPage)
