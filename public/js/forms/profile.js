$(document).ready(function() {

	// Profile form live validation

	// Setup rules
	jQuery.validator.addMethod("isEmailaddress", function(value, element) {
	  return /^[a-z0-9!#$\%&'*+\/=?^_`{|}~.-]+@[a-z0-9-]+(\.[a-z0-9-]+)*$/.test(value);
	}, "Email address is invalid");

	jQuery.validator.addMethod("isUsername", function(value, element) {
	  return /^[a-zA-Z0-9]+$/.test(value);
	}, "Username is invalid");

	jQuery.validator.addMethod("isName", function(value, element) {
	  return /^[^#@]+$/.test(value);
	}, "Name is invalid");

	// Perform validation
	$('#edit-profile').validate({
		rules: {
			"profile.Name": {
				required: true,
				minlength: 6,
				maxlength: 100,
				isName: true
			},
			"profile.UserName": {
				required: true,
				maxlength: 64,
				isUsername: true,
				remote: {
					url: "/account/checkusername",
					type: "post",
					data: {
						"username": function() {
							return $('form input[name="profile.UserName"]').val();
						},
						"currentUsername": function() {
							return $('form input[name="profile.UserName"]').attr('data-current-value');
						}
					}
				}
			},
			"profile.User.Email": {
				required: true,
				minlength: 6,
				maxlength: 200,
				isEmailaddress: true,
				remote: {
					url: "/account/checkemail",
					type: "post",
					data: {
						"email": function() {
							return $('form input[name="profile.User.Email"]').val();
						},
						"currentEmail": function() {
							return $('form input[name="profile.User.Email"]').attr('data-current-value');
						}
					}
				}
			},
			"profile.User.Password": {
				minlength: 6,
				maxlength: 200
			}
		},
		messages: {
			"profile.Name": {
				required: "Name required",
				minlength: "Name must be at least 6 characters",
				maxlength: "Name must be at most 100 characters",
				isName: "Invalid Name. Reserved characters ('#' and '@') are not allowed"
			},
			"profile.UserName": {
				required: "User name required",
				maxlength: "User name can not exceed 64 characters",
				isUsername: "Invalid User name. Alphanumerics allowed only",
				remote: "User name is not available"
			},
			"profile.User.Email": {
				required: "Email address required",
				maxlength: "Email address can not exceed 200 characters",
				minlength: "Email address can not be less than 6 characters",
				isEmailaddress: "You must provide a valid email address",
				remote: "Email address is already registered"
			},
			"profile.User.Password": {
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