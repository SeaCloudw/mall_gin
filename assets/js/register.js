document.getElementById('registerForm').addEventListener('submit', function(event) {
    event.preventDefault(); // 阻止表单默认提交行为

    const name = document.getElementById('name').value;
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;

    // 构造请求体
    const requestBody = {
        name: name,
        email: email,
        password: password
    };
console.log("json:",requestBody)
    // 发送 POST 请求
    // fetch('/api/register', {
    fetch('/api/customers', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(requestBody)
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            showSuccessMessageAndRedirect();
        } else {
            document.getElementById('errorMessage').textContent = data.message || '注册失败，请重试';
        }
    })
    .catch(error => {
        console.error('Error:', error);
        document.getElementById('errorMessage').textContent = '网络错误，请稍后再试';
    });
});

function showSuccessMessageAndRedirect() {
    // 创建一个模态框元素
    const modal = document.createElement('div');
    modal.className = 'success-modal';
    modal.innerHTML = `
        <div class="modal-content">
            <p>注册成功, 5秒后跳转到登录页面！</p>
            <button id="redirect-now">立刻跳转</button>
        </div>
    `;
    document.body.appendChild(modal);

    // 立即跳转按钮事件监听
    document.getElementById('redirect-now').addEventListener('click', function() {
window.location.href='/login'; // 正确：调用 replace 方法并传入目标 URL
    });

    // 设置定时器，在5秒后移除模态框并重定向
    setTimeout(() => {
        document.body.removeChild(modal);
window.location.href='/login'; // 正确：调用 replace 方法并传入目标 URL
    }, 5000); // 5000 毫秒 = 5 秒
}


// // 添加一些基本的样式以美化模态框
// const styleSheet = document.createElement("style");
// styleSheet.type = "text/css";
// styleSheet.innerText = `
//     .success-modal {
//         position: fixed;
//         top: 0;
//         left: 0;
//         width: 100%;
//         height: 100%;
//         background-color: rgba(0, 0, 0, 0.5);
//         display: flex;
//         justify-content: center;
//         align-items: center;
//         z-index: 1000;
//     }

//     .modal-content {
//         background-color: white;
//         padding: 20px;
//         border-radius: 10px;
//         box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
//         text-align: center;
//     }
// `;
// document.head.appendChild(styleSheet);