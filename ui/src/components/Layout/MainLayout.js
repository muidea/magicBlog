import React from 'react'
import styles from './MainLayout.less'
import Header from './Header'
import Footer from './Footer'

function MainLayout({ children, location }) {
  return (
    <div className={styles.normal}>
      <Header location={location} />
      <div className={styles.content}>
        <div className={styles.main}>
          {children}
        </div>
      </div>
      <Footer />
    </div>
  )
}

export default MainLayout
