import React from 'react'
import RichTextEditor, { createEmptyValue } from 'react-rte'
import PropTypes from 'prop-types'
import defaultFormat from './common'
import styles from './view.less'

const RichView = ({ value }) => {
  let val = createEmptyValue()
  value = val.setContentFromString(value, defaultFormat)

  return (<RichTextEditor className={styles.root} value={value} readOnly />)
}

RichView.propTypes = { value: PropTypes.string }

export default RichView
