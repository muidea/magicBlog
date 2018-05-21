import React from 'react'
import { connect } from 'dva'
import styles from './index.less'

function ContactPage() {
  return (
    <div className={styles.normal}>
      <h1 className={styles.title}>ContactPage</h1>
    </div>
  )
}

ContactPage.propTypes = {}

export default connect()(ContactPage)
