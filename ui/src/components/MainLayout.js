import React from 'react'
import Header from './Header'
import Footer from './Footer'
import styles from './MainLayout.less'

function MainLayout({ history, children }) {
  return (
    <div>
      <div className={styles.header}>
        <Header history={history} />
      </div>
      <div className={styles.content}>{children}</div>
      <div className={styles.footer}><Footer /></div>
    </div>
  )
}

export default MainLayout
