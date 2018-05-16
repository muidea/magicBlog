import React from 'react'
import PropTypes from 'prop-types'
import { connect } from 'dva'
import styles from './CatalogPage.css'

function CatalogPage({location, catalog}) {
  const { msg } = catalog

  return (
    <div className={styles.normal}>
      <h1 className={styles.title}>CatalogPage</h1>
      <p>{msg}</p>
    </div>
  )
}

CatalogPage.propTypes = {
  location: PropTypes.object,
  catalog: PropTypes.object,
}

export default connect(({ catalog, location }) => ({ catalog, location }))(CatalogPage)
