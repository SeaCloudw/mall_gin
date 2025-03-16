document.getElementById('loadTableBtn').addEventListener('click', function() {
    // 显示表格容器
    document.getElementById('tableContainer').style.display = 'block';

    // 发送AJAX请求获取数据,后端需要处理这个路由
    fetch('/usersdata')
        .then(response => response.json())
        .then(data => {
            const tableBody = document.querySelector('#usersTable tbody');
            tableBody.innerHTML = ''; // 清空表格内容

            // 遍历数据并生成表格行
            data.forEach(user => {
                const row = document.createElement('tr');

                const idCell = document.createElement('td');
                idCell.textContent = user.id;
                row.appendChild(idCell);

                const nameCell = document.createElement('td');
                nameCell.textContent = user.name;
                row.appendChild(nameCell);

                const emailCell = document.createElement('td');
                emailCell.textContent = user.email;
                row.appendChild(emailCell);

                tableBody.appendChild(row);
            });
        })
        .catch(error => console.error('Error fetching data:', error));
});