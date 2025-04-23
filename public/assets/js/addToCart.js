$(document).ready(function () {
    console.log("jQuery script loaded");

    $(".add_to_cart_button").click(function () {
        let productID = parseInt($("#productID").val());
        let productPrice = parseFloat($("#productPrice").val());
        let productQuantity = parseInt($("#productQuantity").val());

        if (isNaN(productID) || isNaN(productPrice) || isNaN(productQuantity)) {
            alert("Lỗi: Dữ liệu sản phẩm không hợp lệ!");
            return;
        }

        console.log("Adding to cart:", {productID, productPrice, productQuantity});

        $.ajax({
            url: "/auth/add-to-cart",
            method: "POST",
            contentType: "application/json",
            data: JSON.stringify({
                product_id: productID,
                price: productPrice,
                quantity: productQuantity
            }),
            success: function (data) {
                if (data.success) {
                    alert("✅ " + data.message);
                } else if (data.maxQuantity !== undefined) {
                    alert("⚠️ Không thể thêm! Chỉ còn " + data.maxQuantity + " sản phẩm trong kho.");
                } else {
                    alert("❌ " + data.message);
                }
            },
            error: function (xhr) {
                console.error("Lỗi:", xhr.responseText);
                alert("Có lỗi xảy ra, vui lòng thử lại!");
            }
        });
    });
});
