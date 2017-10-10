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

    // Make the functions accessible to other code on the page
    window.ajax = ajax;
    window.error = error;

});
