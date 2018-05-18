import React from 'react'
import PropTypes from 'prop-types'
import Header from './Header'
import styles from './FrameLayout.css'

const FrameLayout = ({ children, location }) => {
  return (
    <div className={styles.normal}>
      <Header location={location} />
      <div className={styles.content}>
        <div className={styles.main}>
          {children}
        </div>
      </div>
    </div>
  )
}

FrameLayout.propTypes = {
  children: PropTypes.object,
  location: PropTypes.object,
}

export default FrameLayout
