import React from 'react'
import { connect } from 'dva'
import { MainLayout } from '../components/Layout'
import styles from './Catalogs.css'

const Catalogs = () => {
  return (
    <MainLayout>
      <div className={styles.normal}>
        Route Component: Catalogs
      </div>
    </MainLayout>
  )
}

function mapStateToProps() {
  return {}
}

export default connect(mapStateToProps)(Catalogs)
