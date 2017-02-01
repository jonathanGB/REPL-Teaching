$(function() {
  let editor = ace.edit("editor");
  editor.setTheme("ace/theme/monokai");
  editor.setHighlightActiveLine(true);
  editor.session.setMode("ace/mode/javascript");
  editor.setShowPrintMargin(false);
  editor.setFontSize(14);

  let results = ace.edit("results");
  results.setReadOnly(true);
  results.setValue("hello true");
  results.setShowPrintMargin(false);
})
