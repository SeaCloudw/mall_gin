// .addEventListener('DOMContentLoaded', function () {
//     console.log("loading begin");
//     loadHeaderAndWelcome(); // 调用外部 JS 文件中的函数
//     validate_get('/protect');
//     loadCartItems(); // 加载购物车内容
// });

// 封装加载购物车内容的函数
function loadCartItems() {
    const token = localStorage.getItem('authToken');
    if (!token) {
        console.error("No authentication token found.");
        return;
    }

    fetch('/api/token/cart', {
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
        populateCartItems(data);
    })
    .catch(error => {
        console.error('There has been a problem with your fetch operation:', error);
    });
}
// 封装填充购物车内容的函数
function populateCartItems(products) {
    const cartItemsDiv = document.getElementById('cart-items');
    cartItemsDiv.innerHTML = ''; // 清空现有内容

    products.forEach(product => {
        // 确保 UnitPrice 存在且为有效数值
        const unitPrice = typeof product.unit_price === 'number' ? product.unit_price : 0;
        const quantity = typeof product.quantity === 'number' ? product.quantity : 1;

        const itemDiv = document.createElement('div');
        itemDiv.className = 'cart-item';
        itemDiv.innerHTML = `
            <div class="checkbox"><input type="checkbox" data-id="${product.product_id}"></div>
            <div class="product-name">${product.product_name || '未知商品名称'}</div>
                        <div class="product-description">${product.product_detail || '无商品说明'}</div>
            <div class="product-price">$${unitPrice.toFixed(2)}</div>

            <div class="product-quantity">
                <label for="quantity-${product.product_id}">数量:</label>
                <input type="number" id="quantity-${product.product_id}" name="quantity" min="1" value="1" style="width: 50px;">
            </div>
        `;
        cartItemsDiv.appendChild(itemDiv);
    });
}

// 预留的提交订单逻辑
function submitOrder() {
    const selectedProducts = [];
    const checkboxes = document.querySelectorAll('.cart-item input[type="checkbox"]:checked');

    checkboxes.forEach(checkbox => {
        const productId = checkbox.getAttribute('data-id');
        selectedProducts.push(productId);
    });

    if (selectedProducts.length === 0) {
        alert('请选择至少一个商品进行下单。');
        return;
    }

    const token = localStorage.getItem('authToken');
    if (!token) {
        console.error("No authentication token found.");
        return;
    }

    fetch('/api/token/orders', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({ productIds: selectedProducts })
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok ' + response.statusText);
        }
        return response.json();
    })
    .then(data => {
        alert('订单提交成功！');
        loadCartItems(); // 刷新购物车内容
    })
    .catch(error => {
        console.error('There has been a problem with your fetch operation:', error);
    });
}