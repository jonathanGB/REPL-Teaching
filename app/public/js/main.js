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

		$('#createGroup').submit(function(e) {
			e.preventDefault()

			if ($('#groupPassword').val() !== $('#rGroupPassword').val()) {
				return toastr.error('Les deux mots de passe sont différents')
			}

			fetch("/groups/", {
				method: "POST",
				credentials: "include",
				body: new FormData(document.getElementById('createGroup'))
			})
			.then(response => response.json())
			.then(payload => {
				if (payload.error) {
					toastr.error(payload.error)
				} else {
					toastr.success('Groupe créé!')

					// TODO: show more info in LIs?
					$('#groupsList').append(
						`<li>
							<a href="/groups/${payload.data}/files"><h3>${$('#groupName').val().trim()}</h3></a>
						</li>`
					)

					$('#createGroupModal').modal('hide')
				}
			})
		})

		$('#joinGroup').submit(function(e) {
			e.preventDefault()
			let gId = $(this).data('groupid')

			fetch(`/groups/${gId}/join`, {
				method: "POST",
				credentials: "include",
				body: new FormData(document.getElementById('joinGroup'))
			})
			.then(response => response.json())
			.then(payload => {
				if (payload.error) {
					toastr.error(payload.error)
				} else {
					location.replace(payload.redirect)
				}
			})
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
