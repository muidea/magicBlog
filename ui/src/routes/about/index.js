import React from 'react'
import PropTypes from 'prop-types'
import { connect } from 'dva'
import { Row } from 'antd'
import { RichView } from '../../components'
import styles from './index.less'

function AboutPage({ about }) {
  const { content } = about

  return (
    <div className={styles.normal}>
      <Row gutter={24}><RichView value={content} /> </Row>
    </div>
  )
}

AboutPage.propTypes = {
  about: PropTypes.object,
}

export default connect(({ about }) => ({ about }))(AboutPage)
