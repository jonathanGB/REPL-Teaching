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
})
