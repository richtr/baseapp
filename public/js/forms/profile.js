$(document).ready(function() {

	// Profile form live validation

	// Setup rules
	jQuery.validator.addMethod("emailaddress", function(value, element) {
	  return /^[a-z0-9!#$\%&'*+\/=?^_`{|}~.-]+@[a-z0-9-]+(\.[a-z0-9-]+)*$/.test(value);
	}, "Email address is invalid");

	jQuery.validator.addMethod("username", function(value, element) {
	  return /^[a-zA-Z0-9]+$/.test(value);
	}, "Username is invalid");

	// Perform validation
	$('#edit-profile').validate({
		rules: {
			"profile.Name": {
				required: true,
				minlength: 6,
				maxlength: 100
			},
			"profile.UserName": {
				required: true,
				maxlength: 64,
				username: true
			},
			"profile.User.Email": {
				required: true,
				minlength: 6,
				maxlength: 200,
				emailaddress: true
			},
			"profile.User.Password": {
				required: true,
				minlength: 6,
				maxlength: 15
			}
		},
		messages: {
			"profile.Name": {
				required: "Name required",
				minlength: "Name must be at least 6 characters",
				maxlength: "Name must be at most 100 characters"
			},
			"profile.UserName": {
				required: "User name required",
				maxlength: "User name can not exceed 64 characters",
				username: "Invalid User name. Alphanumerics allowed only"
			},
			"profile.User.Email": {
				required: "Email address required",
				emailaddress: "You must provide a valid email address",
				maxlength: "Email address can not exceed 200 characters",
				minlength: "Email address can not be less than 6 characters"
			},
			"profile.User.Password": {
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