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
async function searchProducts() {
    const searchText = document.getElementById('search-input').value;
    try {
        // 发送POST请求并传递搜索文本
        const response = await fetch('/api/products/select', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ searchText: searchText })
        });
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        const data = await response.json();
        console.log("Received Data:", data); // 控制台输出接收到的数据

        // const productsContainer = document.getElementById("products");
        // productsContainer.innerHTML = ''; // 清空当前的产品列表

        // 如果有返回的商品数据，则遍历并显示
        if (data && data.length > 0) {
            data.forEach(product => {
                const productItem = document.createElement("div");
                productItem.className = "product-item";
                productItem.innerHTML = `
                    <div>
                        <h4><a href="/products/${product.product_id}" target="_self">${product.product_name}</a></h4>
                        <p>商品说明：${product.product_detail}</p>
                        <p>单价: ${product.unit_price}</p>
                        <p>分类: ${product.Category ? product.Category.CategoryDescription : '未分类'}</p>
                        <p>供应商: ${product.Supplier ? product.Supplier.Name : '未知供应商'}</p>
                    </div>
                    <div>
                        <button onclick="openModal(${product.product_id})">加入购物车</button>
                        <button onclick="openModal(${product.product_id}, true)">立刻下单</button>
                    </div>
                `;
                productsContainer.appendChild(productItem);
            });
        } else {
            // 如果没有找到相关商品，给出提示
            productsContainer.innerHTML = '<p>没有找到符合条件的商品。</p>';
        }

    } catch (error) {
        console.error("There was a problem with the fetch operation:", error); // 错误处理
    }
}

            document.querySelector('.search-button').addEventListener('click', searchProducts);
