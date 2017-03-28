$(function() {
	let wsProtocol = location.protocol.includes("https") ? "wss" : "ws"
	let wsURL =`${wsProtocol}://${location.host}${location.pathname}ws`

	const socket = new WebSocket(wsURL)
	socket.onopen = () => {
		console.log('connected')

		// TODO: send ws message when adding a file

		socket.onmessage = e => {
			let payload = JSON.parse(e.data)

			switch (payload.type) {
				case "live-editing":
					if (payload.data.files) {
						payload.data.files.forEach(fId => {
							$(`#file-${fId} .group`).removeClass('online').addClass(payload.data.status)
						})
					}

					break
				case "update-content":
					if (!payload.err) {
						$(`#file-${payload.data.fId} .file-size`).text(payload.data.size).siblings('.file-lastModified').text(payload.data.lastModified)
					}

					break
				case "update-status":
				console.log('uppppdate')
					if (payload.data.newStatus) { // now private
						$(`#file-${payload.data.fId}`).slideUp()
					} else { // now public
						$(`#file-${payload.data.fId}`).slideDown()

					}
			}

			console.log(payload)
		}

		socket.onclose = (e) => {
			console.log('close', e)
		}
	}
})
