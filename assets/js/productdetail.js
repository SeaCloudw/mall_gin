document.addEventListener('DOMContentLoaded', function() {
    const modal = document.getElementById("modal");
    const closeModal = document.getElementsByClassName("close")[0];
    const addToCartBtn = document.getElementById("add-to-cart-btn");
    const buyNowBtn = document.getElementById("buy-now-btn");

    // Open the modal with product details
    window.openModal = function(productId, isBuyNow = false) {
        modal.style.display = "block";
        fetch(`/api/product/${productId}`).then(response => response.json()).then(product => {
            document.getElementById("modal-title").innerText = product.product_name;
            document.getElementById("modal-unit-price").innerText = product.unit_price.toFixed(2);

            // Store product ID in the modal title element for later use
            document.getElementById("modal-title").dataset.productId = productId;

            if (isBuyNow) {
                addToCartBtn.style.display = 'none';
                buyNowBtn.style.display = 'inline';
            } else {
                addToCartBtn.style.display = 'inline';
                buyNowBtn.style.display = 'none';
            }
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
        const productId = parseInt(document.getElementById("modal-title").dataset.productId);
        const quantity = document.getElementById("quantity").value;
        const address = document.getElementById("address").value;

        // Custom API endpoint for adding to cart
        fetch('/api/cart/add', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ customer_id: 1, product_id: productId, quantity: quantity, address: address })
        }).then(() => modal.style.display = "none");
    };

    // Place an order immediately
    buyNowBtn.onclick = function() {
        const productId = parseInt(document.getElementById("modal-title").dataset.productId);
        const quantity = document.getElementById("quantity").value;
        const address = document.getElementById("address").value;

        // Custom API endpoint for placing an order
        fetch('/api/orders/add', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
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

    // Fetch and display product detail based on URL
    const productId = window.location.pathname.split('/').pop();
    console.log(`Extracted Product ID: ${productId}`);

    function fetchAndDisplayProductDetail(productId) {
        fetch(`/api/products/${productId}`)
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.json();
            })
            .then(data => {
                console.log("Fetched Product Data:", data);

                const productDetailDiv = document.getElementById('product-detail');
                if (!productDetailDiv) {
                    console.error("Element with id 'product-detail' not found.");
                    return;
                }

                // 构建商品详情的 HTML 内容
                productDetailDiv.innerHTML = `
                    <h1 class="product-name">${data.product_name}</h1>
                    <p class="product-detail-info">${data.product_detail}</p>
                    <p class="price">Price: $${data.unit_price.toFixed(2)}</p>
                    <p class="category-description">Category: ${data.Category.CategoryDescription}</p>
                    <div class="buttons">
                        <button class="button add-to-cart" onclick="openModal(${data.product_id})">加入购物车</button>
                        <button class="button buy-now" onclick="openModal(${data.product_id}, true)">立刻下单</button>
                    </div>
                `;
                console.log("Product detail has been updated in the DOM.");
            })
            .catch(error => {
                console.error('Error fetching product details:', error);
            });
    }

    // 调用函数来获取并展示商品详情
    fetchAndDisplayProductDetail(productId);

    // Fetch products from the API and render them on the page (如果你需要的话)
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
                        <h4><a href="/product/${product.product_id}" target="_self">${product.product_name}</a></h4>
                        <p>商品说明：${product.product_detail}</p>
                        <p>单价: ${product.unit_price}</p>
                        <p>分类: ${product.Category.CategoryDescription}</p>
                        <p>供应商: ${product.Supplier.Name}</p>
                    </div>
                    <div>
                        <button onclick="openModal(${product.product_id})">加入购物车</button>
                        <button onclick="openModal(${product.product_id}, true)">立刻下单</button>
                    </div>
                `;
                productsContainer.appendChild(productItem);
            });

            renderPagination(data.pagination); // 假设有一个函数来渲染分页控件
        } catch (error) {
            console.error("There was a problem with the fetch operation:", error); // 错误处理
        }
    }
});