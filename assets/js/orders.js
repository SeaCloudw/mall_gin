
// 封装加载订单内容的函数
function loadOrderItems() {
    const token = localStorage.getItem('authToken');
    if (!token) {
        console.error("No authentication token found.");
        return;
    }

    fetch('/api/token/orders', {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${token}`
        }
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok ' + response.statusText);
        }
        return response.json();
    })
    .then(data => {
        populateOrderItems(data);
    })
    .catch(error => {
        console.error('There has been a problem with your fetch operation:', error);
    });
}

// 封装填充订单内容的函数
function populateOrderItems(orders) {
    console.log("json=", orders);

    const orderItemsDiv = document.getElementById('order-items');
    orderItemsDiv.innerHTML = ''; // 清空现有内容

    orders.data.forEach(order => {
        const product = order.Product;
        const total = order.unit_price * order.amount;

        console.log("json=",order);
        
        const itemDiv = document.createElement('div');
        itemDiv.className = 'order-item';
        itemDiv.innerHTML = `
            <div class="checkbox"><input type="checkbox" data-id="${product.product_id}"></div>
            <div class="product-name">${product.product_name || '未知商品名称'}</div>
            <div class="product-price">$${order.unit_price.toFixed(2)}</div>
            <div class="product-quantity">数量: ${order.amount}</div>
            <div class="product-total">小计: $${total.toFixed(2)}</div>
            <div class="order-status">${order.status || '未知状态'}</div>
        `;
        orderItemsDiv.appendChild(itemDiv);
    });
}