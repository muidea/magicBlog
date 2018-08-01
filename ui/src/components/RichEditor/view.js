import React, { Component } from 'react'
import RichTextEditor from 'react-rte'
import PropTypes from 'prop-types'
import defaultFormat from './common'


class RichView extends Component {
  constructor(props) {
    super(props)
    this.state = { value: this.props.value, format: defaultFormat }
  }

  componentWillReceiveProps(nextProps) {
    this.setState({ value: nextProps.value })
  }

  render() {
    const { value, format } = this.state

    let curValue = RichTextEditor.createEmptyValue()
    if (value) {
      curValue = curValue.setContentFromString(value, format)
    }

    return (
      <RichTextEditor value={curValue} readOnly />
    )
  }
}

RichView.propTypes = { value: PropTypes.string }

export default RichView
