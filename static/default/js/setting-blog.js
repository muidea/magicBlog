$(function() {

  $("#settingForm input,#settingForm textarea").jqBootstrapValidation({
    preventSubmit: true,
    submitError: function($form, event, errors) {
      // additional error messages or events
    },

    submitSuccess: function($form, event) {
      event.preventDefault(); // prevent default submit behaviour
      // get values from FORM
      var id = $("#setting-id").val();
      var name = $("#setting-name").val();
      var domain = $('#setting-domain').val();
      var keyword = $('#setting-keyword').val();
      var email = $('#setting-email').val();
      var icp = $('#setting-icp').val();

      $this = $("#submitBlogButton");
      $this.prop("disabled", true); // Disable submit button until AJAX call is complete to prevent duplicate messages
      $.ajax({
        url: "/api/v1/blog/setting/",
        dataType:'json',
        type: "POST",
        contentType : "application/json",
        data: JSON.stringify({
          id: Number(id),
          name: name,
          domain: domain,
          keyword: keyword,
          email: email,
          icp:icp
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
});

/*When clicking on Full hide fail/success boxes */
$('#settingForm .form-group .form-control').focus(function() {
  $('#success').html('');
});
