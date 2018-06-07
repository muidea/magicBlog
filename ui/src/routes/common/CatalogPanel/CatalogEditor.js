import React from 'react'
import PropTypes from 'prop-types'
import { Button, Form, Input } from 'antd'

const FormItem = Form.Item
const { TextArea } = Input

const CatalogEditor = ({
  form: {
    getFieldDecorator,
    validateFieldsAndScroll,
  },
  onSubmit,
  content,
}) => {
  const handleOk = () => {
    validateFieldsAndScroll((errors, values) => {
      if (errors) {
        return
      }

      const { catalog } = content
      values = { ...values, catalog }

      onSubmit(values)
    })
  }

  return (
    <div>
      <div>
        <form>
          <FormItem label="分类名" hasFeedback>
            {getFieldDecorator('name', {
              initialValue: content.name,
              rules: [{ required: true }],
            })(<Input size="large" onPressEnter={handleOk} placeholder="输入分类名" />)}
          </FormItem>
          <FormItem label="描述">
            {getFieldDecorator('description', {
              initialValue: content.description,
            })(<TextArea rows={3} cols={3} />)}
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

CatalogEditor.propTypes = {
  form: PropTypes.object,
  onSubmit: PropTypes.func,
  content: PropTypes.object,
}

export default Form.create()(CatalogEditor)
