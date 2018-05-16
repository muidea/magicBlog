import React from 'react'
import PropTypes from 'prop-types'
import { connect } from 'dva'
import styles from './index.css'

function CatalogPage({catalog}) {
  const { msg } = catalog

  return (
    <div className={styles.normal}>
      <h1 className={styles.title}>CatalogPage</h1>
      <p>{msg}</p>
    </div>
  )
}

CatalogPage.propTypes = {
  catalog: PropTypes.object,
}

export default connect(({ catalog }) => ({ catalog }))(CatalogPage)
