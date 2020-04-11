$('.post-comment').on('click', function(){
  $.confirm({
      title: '发表留言',
      content: '' +
          '<form action="" class="formName">' +
          '<div class="form-group">' +
          '<input type="text" class="form-control" placeholder="姓名" id="name" required />' +
          '<input type="email" class="form-control" placeholder="电子邮件" id="email" required />' +
          '<textarea rows="5" class="form-control" placeholder="留言信息" id="message" ></textarea>' +
          '</div>' +
          '</form>',
      buttons: {
          formSubmit: {
              text: 'Submit',
              btnClass: 'btn-blue',
              action: function(){
                  var name = this.$content.find('.name').val();
                  if(!name){
                      $.alert('provide a valid name');
                      return false;
                  }
                  $.alert('Your name is ' + name);
              }
          },
          cancel: function(){
              //close
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

$(function() {

  $("#contactForm input,#contactForm textarea").jqBootstrapValidation({
    preventSubmit: true,
    submitError: function($form, event, errors) {
      // additional error messages or events
    },
    submitSuccess: function($form, event) {
      event.preventDefault(); // prevent default submit behaviour
      // get values from FORM
      var host = $("input#host").val();
      var name = $("input#name").val();
      var email = $("input#email").val();
      var message = $("textarea#message").val();
      $this = $("#sendMessageButton");
      $this.prop("disabled", true); // Disable submit button until AJAX call is complete to prevent duplicate messages
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
          $('#success').html("<div class='alert alert-danger'>");
          $('#success > .alert-danger').html("<button type='button' class='close' data-dismiss='alert' aria-hidden='true'>&times;")
            .append("</button>");
          $('#success > .alert-danger').append($("<small>").text(result.reason));
          $('#success > .alert-danger').append('</div>');            
          }
        },
        error: function() {
          // Fail message
          $('#success').html("<div class='alert alert-danger'>");
          $('#success > .alert-danger').html("<button type='button' class='close' data-dismiss='alert' aria-hidden='true'>&times;")
            .append("</button>");
          $('#success > .alert-danger').append($("<small>").text("提交失败!"));
          $('#success > .alert-danger').append('</div>');
        },
        complete: function() {
          setTimeout(function() {
            $this.prop("disabled", false); // Re-enable submit button when AJAX call is complete
          }, 1000);
        }
      });
    },
    filter: function() {
      return $(this).is(":visible");
    },
  });

  $("a[data-toggle=\"tab\"]").click(function(e) {
    e.preventDefault();
    $(this).tab("show");
  });
});

/*When clicking on Full hide fail/success boxes */
$('#name').focus(function() {
  $('#success').html('');
});
