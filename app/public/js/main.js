const ALLOWED_EXTENSIONS = new Set(["go", "js"])


$(function() {
    // at the beginning
    updateFooterOpacity()
    $('[data-toggle="popover"]').popover();


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
							<a href="/groups/${payload.data.id}/files/">
								<div class="group">
									<span class="glyphicon glyphicon-chevron-right group-chevron"></span>
									<h3>${payload.data.groupName}</h3>
									Prof: ${payload.data.teacherName} <br>
								</div>
							</a>
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

		$('#searchGroup').submit(function(e) {
			e.preventDefault()
			let gId = $('#groupId').val()

			location.assign(`/groups/${gId}/join`)
		})

		$('input, select').change(function(e) {
			$(this).addClass('dirty')
		})

		// try to auto-fill fileName and fileExtension if they haven't been modified yet by the user
		$('#fileContentInput').change(function(e) {
			let fileName = $(this)[0].files[0].name
			let fileExtension = fileName.match(/\.(\w+)$/)[1]

			if (!$('#fileNameInput').hasClass('dirty')) {
				$('#fileNameInput').val(fileName)
			}

			if (fileExtension && ALLOWED_EXTENSIONS.has(fileExtension) && !$('#fileExtensionInput').hasClass('dirty')) {
				$('#fileExtensionInput').val(fileExtension)
			}
		})

		$('#createFile').submit(function(e) {
			const MAX_FILE_SIZE = 10000 // 10kB
			e.preventDefault()

			let gId = $(this).data('groupid')

			if ($('#fileContentInput')[0].files[0].size > MAX_FILE_SIZE) {
				toastr.error('Les fichiers ne peuvent pas dépasser 10kB')
				return
			}

			if (!ALLOWED_EXTENSIONS.has($('#fileExtensionInput').val())) {
				toastr.error("Le fichier n'a pas une extension supportée")
				return
			}

			fetch(`/groups/${gId}/files`, {
				method: "POST",
				credentials: "include",
				body: new FormData(document.getElementById('createFile'))
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
