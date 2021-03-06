$(function() {

  $("#editForm input,#editForm textarea").jqBootstrapValidation({
    preventSubmit: true,
    submitError: function($form, event, errors) {
      // additional error messages or events
    },

    submitSuccess: function($form, event) {
      event.preventDefault(); // prevent default submit behaviour
      // get values from FORM
      var id = $("#blog-id").val();
      var title = $("#blog-title").val();
      var content = $('#blog-content').val();
      var catalog = $('#blog-catalog').val();

      $this = $("#submitBlogButton");
      $this.prop("disabled", true); // Disable submit button until AJAX call is complete to prevent duplicate messages
      $.ajax({
        url: "/api/v1/blog/post/",
        dataType:'json',
        type: "POST",
        contentType : "application/json",
        data: JSON.stringify({
          id: Number(id),
          title: title,
          content: content,
          catalog: catalog
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
$('#blog-title').focus(function() {
  $('#success').html('');
});
