$(function() {
    // at the beginning
    updateFooterOpacity()


    // events
    $(window).scroll(updateFooterOpacity)
})


// functions
function updateFooterOpacity() {
    if($(window).scrollTop() + $(window).height() >= $(document).height() - 50) {
        $('footer').css('opacity', 1)
    } else {
        $('footer').css('opacity', 0.7)
    }
}