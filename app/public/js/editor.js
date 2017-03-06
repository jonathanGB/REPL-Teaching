$(function() {
	let editorElem = $('#editor')
  let editor = ace.edit("editor")
	let extensionLangMap = {
		js: 'javascript',
		go: 'golang'
	}
  editor.setTheme("ace/theme/monokai");
  editor.setHighlightActiveLine(true);
  editor.session.setMode(`ace/mode/${extensionLangMap[editorElem.data('extension')]}`);
  editor.setShowPrintMargin(false);
  editor.setFontSize(14);
	editor.setValue(atob(editorElem.data('code')))

  let results = ace.edit("results");
  results.setReadOnly(true);
  results.setShowPrintMargin(false);


	$('#changeStatus').click(function(e) {
		let currStatus = $(this).data('status')
		// update db
		// change style
		$(this).data('status', !currStatus)
		$(this).children('.status-text').text(currStatus ? "Public" : "Privé").siblings('.file-status').removeClass(`${currStatus}`).addClass(`${!currStatus}`)
	})
})
