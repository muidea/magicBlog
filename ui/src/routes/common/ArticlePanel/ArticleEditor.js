import React from 'react'
import PropTypes from 'prop-types'
import { Button, Form, Input } from 'antd'
import { RichEditor } from '../../../components'

const FormItem = Form.Item

const ArticleEditor = ({
  onSubmit,
  content,
  form: {
    getFieldDecorator,
    validateFieldsAndScroll,
  },
}) => {
  const handleOk = () => {
    validateFieldsAndScroll((errors, values) => {
      if (errors) {
        return
      }

      onSubmit(values)
    })
  }

  return (
    <div>
      <div>
        <form>
          <FormItem label="标题" hasFeedback>
            {getFieldDecorator('title', {
              initialValue: content.title,
              rules: [{ required: true }],
            })(<Input size="large" onPressEnter={handleOk} placeholder="输入标题" />)}
          </FormItem>
          <FormItem label="内容">
            {getFieldDecorator('content', {
              initialValue: content.content,
            })(<RichEditor />)}
          </FormItem>
        </form>
      </div>
      <div>
        <Button type="primary" size="large" onClick={handleOk}>
            保存
        </Button>
      </div>
    </div>
  )
}

ArticleEditor.propTypes = {
  form: PropTypes.object,
  onSubmit: PropTypes.func,
  content: PropTypes.object,
}

export default Form.create()(ArticleEditor)
