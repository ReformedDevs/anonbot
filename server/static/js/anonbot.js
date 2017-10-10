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

    // Make the function accessible to other code on the page
    window.ajax = ajax;

});
