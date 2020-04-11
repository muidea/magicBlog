$('.post-comment').on('click', function(){
  $.confirm({
      title: '发表留言',
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
          '  <div class="success"></div>' +
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
                  var host = $(".comment-panel .host").val();

                  $.ajax({
                    url: "/api/v1/comment/post/",
                    dataType:'json',
                    type: "POST",
                    contentType : "application/json",
                    data: JSON.stringify({
                      host: Number(host),
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
                      this.$content.find('.success').html("<div class='alert alert-danger'>");
                      this.$content.find('.success > .alert-danger').html("<button type='button' class='close' data-dismiss='alert' aria-hidden='true'>&times;")
                        .append("</button>");
                        this.$content.find('.success > .alert-danger').append($("<small>").text(result.reason));
                      this.$content.find('.success > .alert-danger').append('</div>');
                      }
                    },
                    error: function() {
                      // Fail message
                      this.$content.find('.success').html("<div class='alert alert-danger'>");
                      this.$content.find('.success > .alert-danger').html("<button type='button' class='close' data-dismiss='alert' aria-hidden='true'>&times;")
                        .append("</button>");
                      this.$content.find('.success > .alert-danger').append($("<small>").text("提交失败!"));
                      this.$content.find('.success > .alert-danger').append('</div>');
                    },
                    complete: function() {
                      setTimeout(function() {
                        $this.prop("disabled", false); // Re-enable submit button when AJAX call is complete
                      }, 1000);
                    }
                  });
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
      }
  });
});

/*When clicking on Full hide fail/success boxes */
$('.commentForm .name').focus(function() {
  $('.commentForm .success').html('');
});
