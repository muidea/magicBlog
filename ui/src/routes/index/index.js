import React from 'react'
import PropTypes from 'prop-types'
import { connect } from 'dva'
import styles from './index.css'

function IndexPage({ index }) {
  const { msg } = index

  return (
    <div className={styles.normal}>
      <h1 className={styles.title}>{ msg }</h1>
      <div className={styles.welcome} />
      <ul className={styles.list}>
        <li>
          To get started, edit <code>src/index.js</code> and save to reload.
        </li>
        <li>
          <a href="https://github.com/dvajs/dva-docs/blob/master/v1/en-us/getting-started.md">
            Getting Started
          </a>
        </li>
      </ul>
    </div>
  )
}

IndexPage.propTypes = {
  index: PropTypes.object,
}

export default connect(({ index }) => ({ index }))(IndexPage)
