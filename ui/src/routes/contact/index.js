import React from 'react'
import PropTypes from 'prop-types'
import { connect } from 'dva'
import { Row } from 'antd'
import { RichView } from '../../components'
import styles from './index.less'

function ContactPage({ contact }) {
  const { content } = contact

  return (
    <div className={styles.normal}>
      <Row gutter={24}><RichView value={content} /> </Row>
    </div>
  )
}

ContactPage.propTypes = {
  contact: PropTypes.object,
}

export default connect(({ contact }) => ({ contact }))(ContactPage)
