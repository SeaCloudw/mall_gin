function loadCustomerInfo() {
    const token = localStorage.getItem('authToken');
    
    if (!token) {
        console.error('No authentication token found.');
        return;
    }

    fetch('/api/token/profile', {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
        }
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok ' + response.statusText);
        }
        return response.json();
    })
    .then(data => {
        console.log("data=",data);
        data=data.data
        // 更新页面上的用户信息
        if (data && typeof data === 'object') {
            document.getElementById('name').innerText = data.name || '';
            document.getElementById('email').innerText = data.email || '';
            document.getElementById('address').innerText = data.address || '';
            document.getElementById('phonenumber').innerText = data.phonenumber || '';

            // 同时更新编辑表单中的默认值
            document.getElementById('editName').value = data.name || '';
            document.getElementById('editEmail').value = data.email || '';
            document.getElementById('editAddress').value = data.address || '';
            document.getElementById('editPhonenumber').value = data.phonenumber || '';
        } else {
            console.error('Invalid data format:', data);
        }
    })
    .catch(error => {
        console.error('There has been a problem with your fetch operation:', error);
    });
}

// 确保在页面加载完成后调用loadCustomerInfo函数
// window.onload = function() {
//     loadCustomerInfo();
// };

// 显示修改信息的弹窗
function editInfo() {
    // 获取当前用户的信息并填充到编辑表单中
    const name = document.getElementById('name').innerText;
    const email = document.getElementById('email').innerText;
    const address = document.getElementById('address').innerText;
    const phonenumber = document.getElementById('phonenumber').innerText;

    document.getElementById('editName').value = name;
    document.getElementById('editEmail').value = email;
    document.getElementById('editAddress').value = address;
    document.getElementById('editPhonenumber').value = phonenumber;

    // 显示编辑信息的模态框
    document.getElementById('editModal').style.display = 'block';
}

// 提交编辑信息
function submitEdit() {
    const name = document.getElementById('editName').value;
    const email = document.getElementById('editEmail').value;
    const address = document.getElementById('editAddress').value;
    const phonenumber = document.getElementById('editPhonenumber').value;

    // 这里可以添加代码来发送更新请求到服务器
    fetch('/api/token/profile', {
        method: 'PUT',
        headers: {
            'Authorization': `Bearer ${localStorage.getItem('authToken')}`,
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            name: name,
            email: email,
            address: address,
            phonenumber: phonenumber
        })
    })
    .then(response => response.json())
    .then(data => {
        // data=data.data
        // 更新页面上的用户信息
        document.getElementById('name').innerText = name;
        document.getElementById('email').innerText = email;
        document.getElementById('address').innerText = address;
        document.getElementById('phonenumber').innerText = phonenumber;

        // 关闭模态框
        closeModal();
    })
    .catch(error => {
        console.error('There has been a problem with your fetch operation:', error);
    });
}

// 显示注销确认的弹窗
function confirmLogout() {
    // 显示注销确认的模态框
    document.getElementById('logoutModal').style.display = 'block';
}

// 注销操作
function logout() {
    // 执行登出操作，如清除session等
    
    localStorage.removeItem('authToken');  // 清除本地存储中的token
    localStorage.removeItem('username') ;
     localStorage.clear();
    window.location.href = "/login";  // 重定向到登录页
}

// 关闭模态框
function closeModal() {
    document.querySelectorAll('.modal').forEach(modal => modal.style.display = 'none');
}