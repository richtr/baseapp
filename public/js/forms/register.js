$(document).ready(function() {

	// Registration form live validation

	// Setup rules
	jQuery.validator.addMethod("emailaddress", function(value, element) {
	  return /^[a-z0-9!#$\%&'*+\/=?^_`{|}~.-]+@[a-z0-9-]+(\.[a-z0-9-]+)*$/.test(value);
	}, "Email address is invalid");

	jQuery.validator.addMethod("username", function(value, element) {
	  return /^[a-zA-Z0-9]+$/.test(value);
	}, "Username is invalid");

	// Perform validation
	$('#register-user').validate({
		rules: {
			"name": {
				required: true,
				minlength: 6,
				maxlength: 100
			},
			"username": {
				required: true,
				maxlength: 64,
				username: true
			},
			"user.Email": {
				required: true,
				minlength: 6,
				maxlength: 200,
				emailaddress: true
			},
			"user.Password": {
				required: true,
				minlength: 6,
				maxlength: 15
			}
		},
		messages: {
			"name": {
				required: "Name required",
				minlength: "Name must be at least 6 characters",
				maxlength: "Name must be at most 100 characters"
			},
			"username": {
				required: "User name required",
				maxlength: "User name can not exceed 64 characters",
				username: "Invalid User name. Alphanumerics allowed only"
			},
			"user.Email": {
				required: "Email address required",
				emailaddress: "You must provide a valid email address",
				maxlength: "Email address can not exceed 200 characters",
				minlength: "Email address can not be less than 6 characters"
			},
			"user.Password": {
				required: "Password required",
				minlength: "Password must be at least 6 characters",
				maxlength: "Password must be at most 15 characters"
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