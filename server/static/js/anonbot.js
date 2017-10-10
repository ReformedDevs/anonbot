/**
 * Custom JavaScript
 */

$(function() {

    // Automatically focus the first input element
    $('input:first').focus();

    // Send an ajax request with the specified data
    function ajax(a, d) {
        return $.ajax({
            type: 'POST',
            url: '/ajax',
            data: JSON.stringify({
                action: a,
                data: d
            }),
            contentType: 'application/json; charset=utf-8'
        });
    }

    // Setup the error dialog
    var $modal = $('.error.modal').modal({
            closable: false,
            duration: 0
        });

    // Show the error message dialog
    function error(m) {
        $modal.modal('show').find('.message').html(m);
    }

    // Create a spinner to show in place of an object while a request is active
    function spinner($e, jqXHR) {
        var $spinner = $('<i>').addClass('spinner loading icon');
        $e.after($spinner).hide();
        jqXHR.always(function() {
            $spinner.remove();
            $e.show();
        });
    }

    // Make the functions accessible to other code on the page
    window.ajax = ajax;
    window.error = error;
    window.spinner = spinner;

});
