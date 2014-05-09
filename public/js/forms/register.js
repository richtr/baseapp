$(document).ready(function() {

	// Registration form live validation

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
	$('#register-user').validate({
		rules: {
			"name": {
				required: true,
				minlength: 6,
				maxlength: 100,
				isName: true
			},
			"user.Email": {
				required: true,
				minlength: 6,
				maxlength: 200,
				isEmailaddress: true,
				remote: {
					url: "/account/checkemail",
					type: "post",
					data: {
						"email": function() {
							return $('form input[name="user.Email"]').val();
						}
					}
				}
			},
			"user.Password": {
				required: true,
				minlength: 6,
				maxlength: 200,
				password: '#user.Password'
			},
			"username": {
				required: true,
				maxlength: 64,
				isUsername: true,
				remote: {
					url: "/account/checkusername",
					type: "post",
					data: {
						"username": function() {
							return $('form input[name="username"]').val();
						}
					}
				}
			}
		},
		messages: {
			"name": {
				required: "Name required",
				minlength: "Name must be at least 6 characters",
				maxlength: "Name must be at most 100 characters",
				isName: "Invalid Name. Reserved characters ('#' and '@') are not allowed"
			},
			"username": {
				required: "User name required",
				maxlength: "User name can not exceed 64 characters",
				isUsername: "Invalid User name. Alphanumerics allowed only",
				remote: "User name is not available"
			},
			"user.Email": {
				required: "Email address required",
				maxlength: "Email address can not exceed 200 characters",
				minlength: "Email address can not be less than 6 characters",
				isEmailaddress: "You must provide a valid email address",
				remote: "Email address is already registered"
			},
			"user.Password": {
				required: "Password required",
				maxlength: "Password must be at most 15 characters",
				minlength: "Password must be at least 6 characters"
			}
	  },
		errorPlacement: function(error, element) {
			$(error).appendTo( $(element).parent() );
			$(error).addClass('control-label');
		},
		highlight: function(element) {
			// Clean up password meter (if it is present on current element)
			var meterBar = $('.password-meter-bar', $(element).parent());
			if(meterBar) meterBar.removeClass().addClass('password-meter-bar');
			var meterMessage = $('.password-meter-message', $(element).parent());
			if(meterMessage) meterMessage.detach();

			$(element).closest('.form-group').removeClass('has-success').addClass('has-error');
		},
		unhighlight: function(element) {
			$(element).closest('.form-group').removeClass('has-error').addClass('has-success');
		}
	});

});