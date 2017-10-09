/**
 * Custom JavaScript
 */

(function() {

    // Automatically focus the first input element
    $(function() {
        $('input:first').focus();
    });

    // For reply buttons, activate the form
    $('a.reply').each(function() {
        var $a = $(this);
        $a.click(function() {
            $a.parent().siblings('.form.hidden').show().find('textarea').focus();
            $a.parent().remove();
        });
    });

})();
