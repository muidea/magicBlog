import React from 'react'
import { connect } from 'dva'
import { FrameLayout } from './../components'
import styles from './IndexPage.css'

function IndexPage({location}) {
  return (
    <FrameLayout location={location}>
      <div className={styles.normal}>
        <h1 className={styles.title}>IndexPage</h1>
      </div>
    </FrameLayout>
  )
}

IndexPage.propTypes = {
}

export default connect()(IndexPage)
