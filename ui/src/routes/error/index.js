import React from 'react'
import PropTypes from 'prop-types'
import { connect } from 'dva'
import { Row } from 'antd'
import { RichView } from '../../components'
import styles from './index.less'

function ErrorPage({ error }) {
  const { content } = error

  return (
    <div className={styles.normal}>
      <Row gutter={24}><RichView value={content} /> </Row>
    </div>
  )
}

ErrorPage.propTypes = {
  error: PropTypes.object,
}

export default connect(({ error }) => ({ error }))(ErrorPage)
