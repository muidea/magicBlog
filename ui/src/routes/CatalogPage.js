import React from 'react'
import PropTypes from 'prop-types'
import { connect } from 'dva'
import styles from './CatalogPage.css'

function CatalogPage({app, catalog}) {
  const { msg } = catalog

  console.log(app)

  return (
    <div className={styles.normal}>
      <h1 className={styles.title}>CatalogPage</h1>
      <p>{msg}</p>
    </div>
  )
}

CatalogPage.propTypes = {
  app: PropTypes.object,
  catalog: PropTypes.object,
}

export default connect(({ app, catalog }) => ({ app, catalog }))(CatalogPage)
