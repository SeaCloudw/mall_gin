document.getElementById('loginForm').addEventListener('submit', function(event) {
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
  console.log('request = ', requestBody);

    // 发送 POST 请求
    fetch('/api/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(requestBody)
    })
    .then(response => response.json())
    .then(data => {
        if (data.status=='success') {
            // 保存 token 到 localStorage
            if (data.token) {
                localStorage.setItem('username', requestBody.name); // 假设response.data.username是从服务器获得的用户名
                localStorage.setItem('authToken', data.token);
                console.log('Token saved:', data.token);
                console.log('Token saved:', requestBody.name);
                    const username = localStorage.getItem('username') || '游客';
                console.log('name saved:', username)

            }
            showSuccessMessageAndRedirect();
        } else {
            document.getElementById('errorMessage').textContent = data.message || '登录失败，请重试';
        }
    })
    .catch(error => {
        console.error('Error:', error);
        document.getElementById('errorMessage').textContent = '网络错误，请稍后再试';
    });
});


function fetchWithAuth(url, options = {}) {
    const token = localStorage.getItem('authToken');
    console.log('Fetching with auth: ', url);
    console.log('Current token in fetchWithAuth:', token);

    if (token) {
        options.headers = {
            ...options.headers,
            'Authorization': `Bearer ${token}`
        };
        console.log('Headers with Authorization:', options.headers);
    } else {
        console.log('No valid token found for request.');
    }

    return fetch(url, options);
}

function showSuccessMessageAndRedirect() {
    // 创建一个模态框元素
    const modal = document.createElement('div');
    modal.className = 'success-modal';
    modal.innerHTML = `
        <div class="modal-content">
            <p>登录成功, 5秒后跳转到个人页面！</p>
            <button id="redirect-now">立刻跳转</button>
        </div>
    `;
    document.body.appendChild(modal);

    // 立即跳转按钮事件监听
    document.getElementById('redirect-now').addEventListener('click', function() {
        // verifyTokenAndRedirect();
window.location.href='/index'; // 正确：调用 replace 方法并传入目标 URL
    });

    // 设置定时器，在5秒后移除模态框并重定向
    setTimeout(() => {
        document.body.removeChild(modal);
window.location.href='/index'; // 正确：调用 replace 方法并传入目标 URL
        // verifyTokenAndRedirect(); // 在定时器触发时也进行验证和重定向
    }, 5000); // 5000 毫秒 = 5 秒
}

function verifyTokenAndRedirect() {
    const token = localStorage.getItem('authToken');
    if (!token) {
        console.log('No token found. Redirecting to login page.');
        window.location.href = '/login'; // 如果没有 Token，重定向到登录页面
        return;
    }

    console.log('Verifying token...');
    fetchWithAuth('/index') // 添加查询参数以区分请求类型
        .then(response => {
            if (response.ok) {
                console.log('Token is valid. Redirecting to index page.');
window.location.href='/index'; // 正确：调用 replace 方法并传入目标 URL
                                // console.log('REPALCE');
            } else {
                console.error('Token verification failed. Redirecting to login page.');
                // window.location.href = '/login'; // 如果验证失败，重定向到登录页面
            }
        })
        .catch(error => {
            console.error('Error verifying token:', error);
            // window.location.href = '/login'; // 如果发生错误，重定向到登录页面
        });
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