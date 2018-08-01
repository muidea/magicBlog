import React, { Component } from 'react'
import RichTextEditor from 'react-rte'
import PropTypes from 'prop-types'
import defaultFormat from './common'
import styles from './view.less'

class RichView extends Component {
  constructor(props) {
    super(props)
    this.state = { value: this.props.value }
  }

  componentWillReceiveProps(nextProps) {
    this.setState({ value: nextProps.value })
  }

  render() {
    const { value } = this.state
    let curValue = RichTextEditor.createEmptyValue()
    if (value) {
      curValue = curValue.setContentFromString(value, defaultFormat)
    }

    return (
      <RichTextEditor value={curValue} className={styles.root} readOnly />
    )
  }
}

RichView.propTypes = { value: PropTypes.string }

export default RichView
