
deleteComment = (url, title, id) => {
  $.confirm({
    title: title,
    content: '' +
        '<form name="sentMessage" class="commentForm" novalidate>' +
        '  <div class="control-group"><small>确认删除该留言?</small></div>' +
        '  <div class="result-panel control-group"></div>' +
        '</form>',
    buttons: {
      formSubmit: {
            text: '提交',
            btnClass: 'btn-blue',
            action: function(){
                var result = true;
                $.ajax({
                  url: url,
                  async:false,
                  dataType:'json',
                  type: "POST",
                  contentType : "application/json",
                  data: JSON.stringify({
                    host: Number(id),
                    origin: window.location.href
                  }),
                  cache: false,
                  success: function(result) {
                    if (result.errorCode===0){
                      window.location.href = result.redirect;
                    } else {
                      // Fail message
                      $('.commentForm .result-panel').html("<div class='alert alert-danger'>");
                      $('.commentForm .result-panel > .alert-danger').html("<button type='button' class='close' data-dismiss='alert' aria-hidden='true'>&times;").append("</button>");
                      $('.commentForm .result-panel > .alert-danger').append($("<small>").text(result.reason));
                      $('.commentForm .result-panel > .alert-danger').append('</div>');
                      result = false;
                    }
                  },
                  error: function() {
                    // Fail message
                    $('.commentForm .result-panel').html("<div class='alert-danger'>");
                    $('.commentForm .result-panel > .alert-danger').html("<button type='button' class='close' data-dismiss='alert' aria-hidden='true'>&times;").append("</button>");
                    $('.commentForm .result-panel > .alert-danger').append($("<small>").text("提交失败!"));
                    $('.commentForm .result-panel > .alert-danger').append('</div>');
                    result = false;
                  }
                });

                return result;
            }
        },
        formCancel: {
          text: '取消',
        },
    },
    onContentReady: function(){
        // you can bind to the form
        var jc = this;
        this.$content.find('form').on('submit', function(e){ // if the user submits the form by pressing enter in the field.
            e.preventDefault();
            jc.$$formSubmit.trigger('click'); // reference the button and click it
        });

        /*When clicking on Full hide fail/success boxes */
        this.$content.find('.commentForm .control-group .controls .form-control').focus(function() {
          $('.commentForm .result-panel').html('');
        });

    }
  });
}

onDeleteComment = (id) => {
  deleteComment("/api/v1/comment/delete/","删除留言",id);
}

replyComment = (url, title, id) => {
  $.confirm({
    title: title,
    content: '' +
        '<form name="sentMessage" class="commentForm" novalidate>' +
        '  <div class="control-group">' +
        '    <div class="form-group floating-label-form-group controls">' + 
        '      <label>回复内容</label>' +
        '      <textarea rows="5" class="message form-control" placeholder="回复内容"></textarea>' +
        '      <p class="help-block text-danger"></p>' +
        '    </div>' +
        '  </div>' +
        '  <div class="result-panel control-group"></div>' +
        '</form>',
    buttons: {
      formSubmit: {
            text: '提交',
            btnClass: 'btn-blue',
            action: function(){
                var message = this.$content.find('.message').val();
                var result = true;
                $.ajax({
                  url: url,
                  async:false,
                  dataType:'json',
                  type: "POST",
                  contentType : "application/json",
                  data: JSON.stringify({
                    host: Number(id),
                    message: message,
                    origin: window.location.href
                  }),
                  cache: false,
                  success: function(result) {
                    if (result.errorCode===0){
                      window.location.href = result.redirect;
                    } else {
                      // Fail message
                      $('.commentForm .result-panel').html("<div class='alert alert-danger'>");
                      $('.commentForm .result-panel > .alert-danger').html("<button type='button' class='close' data-dismiss='alert' aria-hidden='true'>&times;").append("</button>");
                      $('.commentForm .result-panel > .alert-danger').append($("<small>").text(result.reason));
                      $('.commentForm .result-panel > .alert-danger').append('</div>');
                      result = false;
                    }
                  },
                  error: function() {
                    // Fail message
                    $('.commentForm .result-panel').html("<div class='alert-danger'>");
                    $('.commentForm .result-panel > .alert-danger').html("<button type='button' class='close' data-dismiss='alert' aria-hidden='true'>&times;").append("</button>");
                    $('.commentForm .result-panel > .alert-danger').append($("<small>").text("提交失败!"));
                    $('.commentForm .result-panel > .alert-danger').append('</div>');
                    result = false;
                  }
                });

                return result;
            }
        },
        formCancel: {
          text: '取消',
        },
    },
    onContentReady: function(){
        // you can bind to the form
        var jc = this;
        this.$content.find('form').on('submit', function(e){ // if the user submits the form by pressing enter in the field.
            e.preventDefault();
            jc.$$formSubmit.trigger('click'); // reference the button and click it
        });

        /*When clicking on Full hide fail/success boxes */
        this.$content.find('.commentForm .control-group .controls .form-control').focus(function() {
          $('.commentForm .result-panel').html('');
        });

    }
  });
}

onReplyComment = (id) => {
  replyComment("/api/v1/comment/reply/","回复留言", id);
}

postComment = (url, title, id) => {
  $.confirm({
    title: title,
    content: '' +
    '<form name="sentMessage" class="commentForm" novalidate>' +
    '  <div class="control-group">' +
    '    <div class="form-group floating-label-form-group controls">' +
    '      <label>姓名*</label>' +
    '      <input type="text" class="name form-control" placeholder="姓名*" required data-validation-required-message="请输入姓名">' +
    '      <p class="help-block text-danger"></p>' + 
    '    </div>' +
    '  </div>' +
    '  <div class="control-group">' +
    '    <div class="form-group floating-label-form-group controls">' +
    '      <label>电子邮件*</label>' +
    '      <input type="email" class="email form-control" placeholder="电子邮件*" required data-validation-required-message="请输入邮箱">' +
    '      <p class="help-block text-danger"></p>' +
    '    </div>' +
    '  </div>' +
    '  <div class="control-group">' +
    '    <div class="form-group floating-label-form-group controls">' + 
    '      <label>留言信息</label>' +
    '      <textarea rows="5" class="message form-control" placeholder="留言信息"></textarea>' +
    '      <p class="help-block text-danger"></p>' +
    '    </div>' +
    '  </div>' +
    '  <div class="result-panel control-group"></div>' +
    '</form>',
  buttons: {
    formSubmit: {
          text: '提交',
          btnClass: 'btn-blue',
          action: function(){
              var name = this.$content.find('.name').val();
              if(!name){
                  this.$content.find('form .name').focus();
                  return false;
              }
              var email = this.$content.find('.email').val();
              if(!email){
                  this.$content.find('form .email').focus();
                  return false;
              }

              var message = this.$content.find('.message').val();
              var result = true;
              $.ajax({
                url: url,
                async:false,
                dataType:'json',
                type: "POST",
                contentType : "application/json",
                data: JSON.stringify({
                  host: Number(id),
                  name: name,
                  email: email,
                  message: message,
                  origin: window.location.href
                }),
                cache: false,
                success: function(result) {
                  if (result.errorCode===0){
                    window.location.href = result.redirect;
                  } else {
                    // Fail message
                    $('.commentForm .result-panel').html("<div class='alert alert-danger'>");
                    $('.commentForm .result-panel > .alert-danger').html("<button type='button' class='close' data-dismiss='alert' aria-hidden='true'>&times;").append("</button>");
                    $('.commentForm .result-panel > .alert-danger').append($("<small>").text(result.reason));
                    $('.commentForm .result-panel > .alert-danger').append('</div>');
                    result = false;
                  }
                },
                error: function() {
                  // Fail message
                  $('.commentForm .result-panel').html("<div class='alert-danger'>");
                  $('.commentForm .result-panel > .alert-danger').html("<button type='button' class='close' data-dismiss='alert' aria-hidden='true'>&times;").append("</button>");
                  $('.commentForm .result-panel > .alert-danger').append($("<small>").text("提交失败!"));
                  $('.commentForm .result-panel > .alert-danger').append('</div>');
                  result = false;
                }
              });

              return result;
          }
      },
      cancel: {
        text: '取消',
      },
  },
  onContentReady: function(){
      // you can bind to the form
      var jc = this;
      this.$content.find('form').on('submit', function(e){ // if the user submits the form by pressing enter in the field.
          e.preventDefault();
          jc.$$formSubmit.trigger('click'); // reference the button and click it
      });

      /*When clicking on Full hide fail/success boxes */
      this.$content.find('.commentForm .control-group .controls .form-control').focus(function() {
        $('.commentForm .result-panel').html('');
      });
    }
  });
}

onPostComment = (id) => {
  postComment("/api/v1/comment/post/","发表留言", id);
}
