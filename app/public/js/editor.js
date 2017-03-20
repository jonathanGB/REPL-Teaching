$(function() {
	$('[data-toggle="popover"]').popover();

	const MAX_FILE_SIZE = 10000 // 10kB
	let editorElem = $('#editor')
	let resultsElem = $('#results')
  let editor = ace.edit("editor")
	let extensionLangMap = {
		js: 'javascript',
		go: 'golang'
	}
	let wsURL =`ws://${location.host}${location.pathname}ws`
	let isOwner = editorElem.data('isowner')

  editor.setTheme("ace/theme/monokai");
  editor.setHighlightActiveLine(true);
  editor.session.setMode(`ace/mode/${extensionLangMap[editorElem.data('extension')]}`);
  editor.setShowPrintMargin(false);
  editor.setFontSize(14);
	editor.setValue(atob(editorElem.data('code')), -1)
	editor.setReadOnly(!isOwner)

  let results = ace.edit("results");
  results.setReadOnly(true);
  results.setShowPrintMargin(false);
	results.setOptions({
    maxLines: 30
	});

	const socket = new WebSocket(wsURL)
	socket.onopen = () => {
		console.log('connected')

		$('#runFile').click(function(e) {
			wsFeedback(this, 'sending')

			let content = editor.getValue()
			if (content.length > MAX_FILE_SIZE) {
				toastr.error("Le fichier dépasse la limite de 10kB")
				return
			}

			let toSend = {
				type: "run",
				content
			}
			socket.send(JSON.stringify(toSend))
		})

		$('#saveFile').click(function(e) {
			wsFeedback(this, 'sending')

			let content = editor.getValue()
			if (content.length > MAX_FILE_SIZE) {
				toastr.error("Le fichier dépasse la limite de 10kB")
				return
			}

			let toSend = {
				type: "update-content",
				content
			}
			socket.send(JSON.stringify(toSend))
		})

		$('#changeStatus').click(function(e) {
			let status = $(this).data('status')

			let toSend = {
				type: "update-status",
				newStatus: !status
			}
			socket.send(JSON.stringify(toSend))
		})

		socket.onmessage = (e) => {
			let payload = JSON.parse(e.data)

			switch (payload.type) {
				case "run":
					wsFeedback(document.getElementById("runFile"), "receiving")
					results.setValue(payload.data, -1)

					if (payload.err) {
						resultsElem.removeClass('no-error').addClass('error')
						toastr.error("Erreur lors de l'exécution du script")
					} else {
						resultsElem.removeClass('error').addClass('no-error')
						toastr.success("Fichier exécuté sans problème!")
					}
					break
				case "update-content":
					wsFeedback(document.getElementById("saveFile"), "receiving")

					if (payload.err) {
						toastr.error("Erreur lors de la sauvegarde du fichier")
					} else {
						toastr.success("Fichier sauvegardé!")
					}
					break
				case "update-status":
					let that = $('#changeStatus')
					let newStatus = payload.data === "true"

					if (payload.err) {
						toastr.error("Erreur lors du changement de statut")
					} else {
						that.data('status', newStatus)
						that.children('.status-text').text(!newStatus ? "Public" : "Privé").siblings('.file-status').removeClass(`${!newStatus}`).addClass(`${newStatus}`)
					}
			}
			console.log(payload)
		}

		socket.onclose = (e) => {
			console.log('close', e)
		}
	}


	// TODO: set on change event
	editor.getSession().on('change', function(e) {
		console.log(e)
	})

	$(document).keyup(function(e) { // if escape key pressed inside iframe
		if (window.parent && e.keyCode == 27) {
			window.parent.removeLightbox(e)
		}
	})
})


function wsFeedback(elem, eventType) {
	if (eventType === "sending") {
		$(elem).addClass('disabled')
		$('#loader').fadeIn(400)
	} else {
		$(elem).removeClass('disabled')
		$('#loader').fadeOut(400)
	}
}
