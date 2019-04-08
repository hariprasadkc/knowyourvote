$(document).ready(function() {
    $("#pincodeSubmit").click(function() {
        $("#pincodeResult").empty();
        var pincode = $("#pincodeInput").val();

        $.ajax({
            url: 'http://knowyourvote.appspot.com/findconstituency?pincode=' + pincode,
            type: 'get',
            success: function(data) {
                $("#pincodeResult").html(data);
            }
        });
        });
});
