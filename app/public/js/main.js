$(function() {
    // at the beginning
    updateFooterOpacity()


    // events
    $(window).scroll(updateFooterOpacity)

		$('#createUser').submit(function(e) {
			if ($('#passwordInput').val() !== $('#repeatPasswordInput').val()) {
				toastr.error('Les deux mots de passe sont différents')
				e.preventDefault()
			}
		})
})


// functions
function updateFooterOpacity() {
    if($(window).scrollTop() + $(window).height() >= $(document).height() - 50) {
        $('footer').css('opacity', 1)
    } else {
        $('footer').css('opacity', 0.7)
    }
}

toastr.options = {
		"closeButton": true,
		"debug": false,
		"newestOnTop": false,
		"progressBar": false,
		"positionClass": "toast-top-right",
		"preventDuplicates": false,
		"onclick": null,
		"showDuration": "300",
		"hideDuration": "1000",
		"timeOut": "5000",
		"extendedTimeOut": "1000",
		"showEasing": "swing",
		"hideEasing": "linear",
		"showMethod": "fadeIn",
		"hideMethod": "fadeOut"
	}
