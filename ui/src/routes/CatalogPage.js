import React from 'react'
import { connect } from 'dva'
import { FrameLayout } from './../components'
import styles from './CatalogPage.css'

function CatalogPage({location}) {
  return (
    <FrameLayout location={location}>
      <div className={styles.normal}>
        <h1 className={styles.title}>CatalogPage</h1>
      </div>
    </FrameLayout>
  )
}

CatalogPage.propTypes = {
}

export default connect()(CatalogPage)
