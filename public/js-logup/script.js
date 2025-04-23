document.addEventListener("DOMContentLoaded", function () {
    const signUpError = document.getElementById('signup-error');
    const signUpButton = document.getElementById('signUp');
    const signInButton = document.getElementById('signIn');
    const container = document.getElementById('container');

    // Kiểm tra lỗi và hiển thị thông báo
    if (signUpError && signUpError.innerText.trim() !== "") {
        signUpError.style.display = "block";
    }

    // Chuyển sang form đăng ký
    signUpButton.addEventListener('click', () => {
        container.classList.add('right-panel-active');
    });

    // Chuyển về form đăng nhập và ẩn thông báo lỗi
    signInButton.addEventListener('click', () => {
        container.classList.remove('right-panel-active');
        if (signUpError) {
            signUpError.style.display = "none"; // Ẩn lỗi khi chuyển về đăng nhập
        }
    });
});
