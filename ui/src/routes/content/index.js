import React from 'react'
import PropTypes from 'prop-types'
import { connect } from 'dva'
import styles from './index.css'

function ContentPage({ content }) {
  const { msg } = content

  return (
    <div className={styles.normal}>
      <h1 className={styles.title}>ContentPage</h1>
      <p>{msg}</p>
    </div>
  )
}

ContentPage.propTypes = {
  content: PropTypes.object,
}

export default connect(({ content }) => ({ content }))(ContentPage)
