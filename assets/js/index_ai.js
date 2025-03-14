document.addEventListener('DOMContentLoaded', function() {
    const modal = document.getElementById("modal");
    const closeModal = document.getElementsByClassName("close")[0];
    const addToCartBtn = document.getElementById("add-to-cart-btn");
    const buyNowBtn = document.getElementById("buy-now-btn");

    // Fetch products from the API and render them on the page
    async function fetchProducts(page = 1) {
        try {
            const response = await fetch(`/api/products?page=${page}`);
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            const data = await response.json();
            console.log("Received Data:", data); // 控制台输出接收到的数据

            const productsContainer = document.getElementById("products");
            productsContainer.innerHTML = ''; // 清空当前的产品列表

            data.data.forEach(product => {
                const productItem = document.createElement("div");
                productItem.className = "product-item";
                productItem.innerHTML = `
                    <div>
                        <h4><a href="/products/${product.product_id}" target="_self">${product.product_name}</a></h4>
                        <p>商品说明：${product.product_detail}</p>
                        <p>单价: ${product.unit_price}</p>
                        <p>分类: ${product.Category.CategoryDescription}</p>
                        <p>供应商: ${product.Supplier.Name}</p>
                    </div>
                    <div>
                        <button id="add-to-cart-btn" data-product-id="${product.product_id}">加入购物车</button>
                        <button id="buy-now-btn" data-product-id="${product.product_id}" data-is-buy-now="true">立刻下单</button>
                    </div>
                `;
                productsContainer.appendChild(productItem);
            });

            // renderPagination(data.pagination); // 假设有一个函数来渲染分页控件
        } catch (error) {
            console.error("There was a problem with the fetch operation:", error); // 错误处理
        }
    }

    // 绑定“加入购物车”和“立刻下单”按钮的点击事件
    document.querySelectorAll('.add-to-cart-btn, .buy-now-btn').forEach(button => {
        button.addEventListener('click', function() {
            const productId = parseInt(this.getAttribute('data-product-id'));
            const isBuyNow = this.getAttribute('data-is-buy-now') === 'true';
            openModal(productId, isBuyNow);
        });
    });

    // Open the modal with product details
    window.openModal = function(productId, isBuyNow = false) {
        modal.style.display = "block";
        fetch(`/api/product/${productId}`).then(response => response.json()).then(product => {
            document.getElementById("modal-title").innerText = product.product_name;
            document.getElementById("modal-unit-price").innerText = product.unit_price;

            // Store product ID in a hidden input element for later use
            document.getElementById("product-id").value = productId;
            console.log("store=", productId);

            if (isBuyNow) {
                addToCartBtn.style.display = 'none';
                buyNowBtn.style.display = 'inline';
            } else {
                addToCartBtn.style.display = 'inline';
                buyNowBtn.style.display = 'none';
            }
        }).catch(error => {
            console.error("There was a problem fetching the product details:", error); // 错误处理
        });
    };

    // Close the modal
    closeModal.onclick = function() {
        modal.style.display = "none";
    };

    // Close the modal when clicking outside of it
    window.onclick = function(event) {
        if (event.target == modal) {
            modal.style.display = "none";
        }
    };

    // Add product to cart
    addToCartBtn.onclick = function() {
        const productId = parseInt(document.getElementById("product-id").value); // 从隐藏的 input 元素中获取 product_id
        const quantity = document.getElementById("quantity").value;
        const address = document.getElementById("address").value;
        const token = localStorage.getItem('authToken'); // 获取存储在localStorage中的token
        console.log("id=", productId);

        // Custom API endpoint for adding to cart
        fetch('/api/cart/add', { // 修正API路径
            method: 'POST',
            headers: { 
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}` // 在请求头中添加Authorization字段
            },
            body: JSON.stringify({
                customer_id: 1, // 假设顾客ID固定为1，可以根据实际情况调整
                product_id: productId,
                quantity: quantity,
                address: address
            })
        }).then(response => {
            if (response.ok) {
                return response.json(); // 如果需要处理返回的数据，可以在这里进行
            } else {
                throw new Error('Failed to add item to cart');
            }
        }).then(data => {
            console.log("Item added to cart successfully:", data); // 成功响应的处理
            modal.style.display = "none"; // 关闭模态框
        }).catch(error => {
            console.error("There was a problem with the fetch operation:", error); // 错误处理
        });
    };

    // Place an order immediately
    buyNowBtn.onclick = function() {
        const productId = parseInt(document.getElementById("product-id").value); // 从隐藏的 input 元素中获取 product_id
        const quantity = document.getElementById("quantity").value;
        const address = document.getElementById("address").value;

        // Custom API endpoint for placing an order
        fetch('/api/orders/add', {
            method: 'POST',
            headers: { 
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('authToken')}`
            },
            body: JSON.stringify({
                customer_id: 1,
                order_date: new Date().toISOString(),
                product_id: [productId],
                unit_price: [parseFloat(document.getElementById("modal-unit-price").innerText)],
                amount: [parseInt(quantity)],
                status: 'Pending',
                address: address
            })
        }).then(() => modal.style.display = "none");
    };

    // 更新欢迎信息
    function updateWelcomeMessage() {
        var welcomeMessageElement = document.querySelector('.welcome-message');
        var username = localStorage.getItem('username') || '游客';
        welcomeMessageElement.textContent = '欢迎, ' + username;
    }

    // 测试验证函数
    function validate_get(url) {
        const token = localStorage.getItem('authToken');
        if (!token) {
            console.error('No authentication token found.');
            window.location.href='/login';//没有token直接跳转到login
            return;
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
            console.log('Success:', data);
        })
        .catch(error => {
            console.error('There has been a problem with your fetch operation:', error);
        });
    }

    // 搜索产品
    async function searchProducts() {
        const searchText = document.getElementById('search-input').value;
        try {
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

            const productsContainer = document.getElementById("products");
            productsContainer.innerHTML = ''; // 清空当前的产品列表

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
                            <button class="add-to-cart-btn" data-product-id="${product.product_id}">加入购物车</button>
                            <button class="buy-now-btn" data-product-id="${product.product_id}" data-is-buy-now="true">立刻下单</button>
                        </div>
                    `;
                    productsContainer.appendChild(productItem);
                });
            } else {
                productsContainer.innerHTML = '<p>没有找到符合条件的商品。</p>';
            }
        } catch (error) {
            console.error("There was a problem with the fetch operation:", error); // 错误处理
        }
    }

    // 绑定搜索按钮点击事件
    document.querySelector('.search-button').addEventListener('click', searchProducts);

    // 初始化页面
    updateWelcomeMessage();
    fetchProducts();

    // 测试访问保护资源
    validate_get('/protect');
});