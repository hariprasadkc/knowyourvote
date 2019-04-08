$(document).ready(function() {
    $("#pincodeSubmit").click(function() {
        
        var pincode = $("#pincodeInput").val();

        $.ajax({
            url: 'http://knowyourvote.appspot.com/findconstituency?pincode=' + pincode,
            type: 'get',
            success: function(data) {
                console.log(data);
            }
        });
        });
});
