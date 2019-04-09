$(document).ready(function() {
    $("#pincodeSubmit").click(function() {
        $("#pincodeResult").empty();
        var pincode = $("#pincodeInput").val();
        var pinlength = $("#pincodeInput").val().length;

        if (pinlength == 6) {
            $.ajax({
                url: 'http://localhost:8080/findconstituency?pincode=' + pincode,
                type: 'get',
                success: function(data) {
                    $("#pincodeResult").html(data);
                }
            });
        } else {
            var errorMsg = "<div class=\"alert alert-danger mx-auto px-5\" role=\"alert\">Enter a valid pincode</div>"
            $("#pincodeResult").html(errorMsg);
        }
    });
});
