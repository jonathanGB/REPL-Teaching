$(function() {
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

	// TODO: set read-only if not owner of file

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

			// TODO: put change inside a callback
			$(this).data('status', !status)
			$(this).children('.status-text').text(status ? "Public" : "Privé").siblings('.file-status').removeClass(`${status}`).addClass(`${!status}`)
		})

		socket.onmessage = (e) => {
			let payload = JSON.parse(e.data)

			switch (payload.type) {
				case "run":
					results.setValue(payload.data, -1)

					if (payload.err) {
						resultsElem.removeClass('no-error').addClass('error')
					} else {
						resultsElem.removeClass('error').addClass('no-error')
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
})
