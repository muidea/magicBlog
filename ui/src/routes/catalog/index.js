import React from 'react'
import { connect } from 'dva'
import styles from './index.less'

function CatalogPage() {
  return (
    <div className={styles.normal}>
      <h1 className={styles.title}>CatalogPage</h1>
    </div>
  )
}

CatalogPage.propTypes = {}

export default connect()(CatalogPage)
