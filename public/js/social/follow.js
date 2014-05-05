$(document).ready(function() {

	function doUpdate(element, successCallback, errorCallback) {
		$.getJSON(element.attr('data-href'), function(json) {
			if(json.Status == 'success') {
				if (successCallback) {
					successCallback(json);
				}
			} else {
				// An error occured
				if (errorCallback) {
					errorCallback(json);
				}
			}
		});
	}

	// Like button functionality
	$('button.like').each(function() {
		$(this).on('click', function( event ) {
			event.preventDefault();
			event.stopPropagation();
			doUpdate( $(this), function() {
				var targetBadge = $('.badge', $(this));
				var count = targetBadge.text() * 1; // force int
				targetBadge.text(count + 1);
				if(count == 0) {
					var icon = $('b', $(this));
					icon.removeClass('glyphicon-heart-empty').addClass('glyphicon-heart');
				}
			}.bind( this ));
		}.bind( this ));
	});

	// Follow button functionality
	$('button.follow').each(function() {
		$(this).on('click', function( event ) {
			event.preventDefault();
			event.stopPropagation();
			doUpdate( $(this), function() {
				var targetBadge = $('#followers');
				var count = targetBadge.text() * 1; // force int
				var btnText = $('.text', $(this));
				if(btnText.text() == 'Follow') {
					targetBadge.text(count + 1);
					$(this).removeClass('btn-primary').addClass('btn-danger');
					btnText.text('Unfollow');
				} else {
					targetBadge.text(count - 1);
					$(this).removeClass('btn-danger').addClass('btn-primary');
					btnText.text('Follow');
				}
			}.bind( this ));
		}.bind( this ));
	});

});