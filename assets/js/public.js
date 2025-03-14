// public.js

// 更新欢迎信息的函数
function updateWelcomeMessage() {
    var welcomeMessageElement = document.querySelector('.welcome-message');
    if (welcomeMessageElement) {
        var username = localStorage.getItem('username') || '游客';
        welcomeMessageElement.textContent = '欢迎, ' + username;
    } else {
        console.error("Element with class 'welcome-message' not found.");
    }
}

// 动态加载 header.html 并更新欢迎信息
function loadHeaderAndWelcome() {
    fetch('/bar')
        .then(response => response.text())
        .then(data => {
            document.getElementById('header-container').innerHTML = data;
            // 在 header.html 加载并插入后调用 updateWelcomeMessage,加载后再用函数,这里再添加跳转
            updateWelcomeMessage();
            console.log("loading over");


            
        })
        .catch(error => console.error('Error loading header:', error));
}

	//专门用于验证token，成功返回一个string,失败直接跳转到login
	// validate_get('/protect');

function validate_get(url) {
    const token = localStorage.getItem('authToken');
    
    if (!token) {
        console.error('No authentication token found.');
window.location.href='/login';
        return
    }

    fetch(url, {
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
        // 处理成功响应的数据
        console.log('Success:', data);
    })
    .catch(error => {
        console.error('There has been a problem with your fetch operation:', error);
    });
}
