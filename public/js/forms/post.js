$(document).ready(function() {

	// Post form live validation

	// Perform validation
	$('#user-post').validate({
		rules: {
			"post.Title": {
				required: true,
				minlength: 3,
				maxlength: 200
			}
		},
		messages: {
			"post.Title": {
				required: "A post must have a title",
				minlength: "Post title must exceed 2 characters",
				maxlength: "Post title cannot exceed 200 characters"
			}
	  },
		errorPlacement: function(error, element) {
			$(error).appendTo( $(element).parent() );
			$(error).addClass('control-label');
		},
		highlight: function(element) {
			$(element).closest('.form-group').removeClass('has-success').addClass('has-error');
		},
		unhighlight: function(element) {
			$(element).closest('.form-group').removeClass('has-error').addClass('has-success');
		}
	});

});