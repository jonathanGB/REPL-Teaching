$(function() {
	const MAX_FILE_SIZE = 10000 // 10kB
	let editorElem = $('#editor')
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


	$('#changeStatus').click(function(e) {
		let currStatus = $(this).data('status')
		// TODO: update db
		// change style
		$(this).data('status', !currStatus)
		$(this).children('.status-text').text(currStatus ? "Public" : "Privé").siblings('.file-status').removeClass(`${currStatus}`).addClass(`${!currStatus}`)
	})

	if (isOwner) {
		const socket = new WebSocket(wsURL)
		socket.onopen = () => {
			console.log('connected')

			$('#saveFile').click(function(e) {
				let data = editor.getValue()
				if (data.length > MAX_FILE_SIZE) {
					toastr.error("Le fichier dépasse la limite de 10kB")
					return
				}

				let toSend = {
					type: "update-file",
					data
				}
				socket.send(JSON.stringify(toSend))
			})

			socket.onmessage = (e) => {
				console.log('message', e)
				editor.setValue(JSON.parse(e.data).data.substring(0, 500))
			}

			socket.onclose = (e) => {
				console.log('close', e)
			}
		}
	}


	// TODO: set on change event
	editor.getSession().on('change', function(e) {
		console.log(e)
	})
})
