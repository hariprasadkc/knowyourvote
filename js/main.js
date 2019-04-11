$(document).ready(function() {
    $("#pincodeSubmit").click(function() {
        $("#pincodeResult").empty();
        var pincode = $("#pincodeInput").val();
        var pinlength = $("#pincodeInput").val().length;

        if (pinlength == 6) {
            $.ajax({
                url: '//knowyourvote.appspot.com/findconstituency?pincode=' + pincode,
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
    
    var checkitem = function() {
        var $this;
        $this = $("#candidateCarousel");
        if ($("#candidateCarousel .carousel-inner .carousel-item:first").hasClass("active")) {
            $this.children(".carousel-control-prev").hide();
            $this.children(".carousel-control-next").show();
        } else if ($("#candidateCarousel .carousel-inner .carousel-item:last").hasClass("active")) {
            $this.children(".carousel-control-next").hide();
            $this.children(".carousel-control-prev").show();
        } else {
            $this.children(".carousel-control-prev").show();
            $this.children(".carousel-control-next").show();
        }
    };
    checkitem();
    $("#candidateCarousel").on("slid.bs.carousel", checkitem);
    
});
