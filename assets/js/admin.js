document.addEventListener('DOMContentLoaded', function () {
    const userTableBody = document.getElementById('user-table-body');
    const addUserBtn = document.getElementById('add-user-btn');
    const deleteSelectedBtn = document.getElementById('delete-selected-btn');
    const selectAllCheckbox = document.getElementById('select-all');
    const modal = document.getElementById('user-form-modal');
    const closeModalBtn = modal.querySelector('.close');
    const userForm = document.getElementById('user-form');
    const modalTitle = document.getElementById('modal-title');
    const userIdInput = document.getElementById('user-id');
    const usernameInput = document.getElementById('username');
    const emailInput = document.getElementById('email');
    const passwordInput = document.getElementById('password');

    let users = []; // 存储用户数据

    // 加载用户数据

function loadUsers() {
    fetch('/api/users')
        .then(response => {
            if (!response.ok) {
                throw new Error('网络响应失败');
            }
            return response.text(); // 先以文本形式获取响应
        })
        .then(text => {
            // console.log('原始响应数据:', text); // 打印原始响应
            return JSON.parse(text); // 尝试解析 JSON
        })
        .then(data => {
            users = data;
            renderTable();
        })
        .catch(error => {
            console.error('加载用户数据失败:', error);
        });
}

    // 渲染表格
    function renderTable() {
        userTableBody.innerHTML = '';
        users.forEach(user => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td><input type="checkbox" class="user-checkbox" data-id="${user.id}"></td>
                <td>${user.username}</td>
                <td>${user.email}</td>
                <td>${new Date(user.created_at).toLocaleDateString()}</td>
                <td><span class="status ${user.active ? 'active' : 'inactive'}">${user.active ? '启用' : '禁用'}</span></td>
                <td>
                    <button class="edit-btn" data-id="${user.id}">编辑</button>
                    <button class="delete-btn" data-id="${user.id}">删除</button>
                </td>
            `;
            userTableBody.appendChild(row);
        });

        // 绑定编辑和删除按钮事件
        document.querySelectorAll('.edit-btn').forEach(btn => {
            btn.addEventListener('click', handleEdit);
        });
        document.querySelectorAll('.delete-btn').forEach(btn => {
            btn.addEventListener('click', handleDelete);
        });
    }

    // 打开添加用户表单
    addUserBtn.addEventListener('click', () => {
        modalTitle.textContent = '添加用户';
        userIdInput.value = '';
        usernameInput.value = '';
        emailInput.value = '';
        passwordInput.value = '';
        modal.style.display = 'flex';
    });

    // 关闭模态框
    closeModalBtn.addEventListener('click', () => {
        modal.style.display = 'none';
    });

    // 提交表单（添加/编辑用户）
    userForm.addEventListener('submit', function (e) {
        e.preventDefault();
        const userData = {
            id: userIdInput.value,
            username: usernameInput.value,
            email: emailInput.value,
        };
            // 如果密码输入框有值，则添加到 userData 中
    if (passwordInput.value) {
        userData.password = passwordInput.value;
    }
console.log(userData);

        //若存在用户，则进行  更新“PUT”
        //不存在 就新增用户 “POST”
        const url = userData.id ? `/api/update/users/${userData.id}` : '/api/register';
        const method = userData.id ? 'PUT' : 'POST';

        fetch(url, {
            method: method,
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(userData),
        })
            .then(response => {
                if (!response.ok) {
                    return response.json().then(errData => {
                        throw new Error(errData.error || '未知错误');
                    });
                }
                return response.json();
            })
            .then(() => {
                loadUsers(); // 重新加载用户数据
                modal.style.display = 'none';
            })
            .catch(error => {
                console.error('保存用户失败:', error.message);
                showErrorNotification(error.message); // 显示错误通知
            });
    });
function showErrorNotification(message) {
    const notification = document.getElementById('errorNotification');
    notification.textContent = message;
    notification.style.display = 'block';

    // 5秒后隐藏弹窗
    setTimeout(() => {
        notification.style.display = 'none';
    }, 5000);
}
    // 编辑用户,点击按钮自动填充内容
    function handleEdit(e) {
        const userId = e.target.getAttribute('data-id');
        // 这个users是最开始返回的保存的值
        const user = users.find(u => u.id == userId);
        if (user) {
            modalTitle.textContent = '编辑用户';
            userIdInput.value = user.id;
            usernameInput.value = user.username;
            emailInput.value = user.email;
        //             // 将密码输入框类型改为 text 以显示密码
        // passwordInput.type = 'text';
            passwordInput.value =''; // 编辑时不显示密码
            modal.style.display = 'flex';
        }
    }

    // 删除用户
    function handleDelete(e) {
        const userId = e.target.getAttribute('data-id');
        // console.log("data=",userId);
if (confirm('确定删除该用户吗？')) {
    fetch(`/api/delete/users/${userId}`, {
        method: 'DELETE'
    })
    .then(response => {
        if (response.ok) {
            // 如果状态码是 2xx（包括 204 No Content），则认为删除成功
            return response.status === 204 ? Promise.resolve() : response.json();
        } else {
            // 处理其他非 2xx 状态码
            return response.json().then(errData => {
                throw new Error(errData.error || '删除用户失败');
            });
        }
    })
    .then(() => {
        // 成功删除后，重新加载用户数据
        loadUsers();
        console.log('删除用户成功');
    })
    .catch(error => {
        // 错误处理
        console.error('删除用户失败:', error);
        alert(`删除用户失败: ${error.message}`);
    });
}
    }

    // 批量删除
    deleteSelectedBtn.addEventListener('click', () => {
        const selectedIds = Array.from(document.querySelectorAll('.user-checkbox:checked'))
            .map(checkbox => checkbox.getAttribute('data-id'));
        if (selectedIds.length > 0 && confirm('确定删除选中的用户吗？')) {
            fetch('/api/users/batch-delete', {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ ids: selectedIds }),
            })
                .then(response => response.json())
                .then(data => {
                    loadUsers(); // 重新加载用户数据
                })
                .catch(error => {
                    console.error('批量删除失败:', error);
                });
        }
    });

    // 全选/取消全选
    selectAllCheckbox.addEventListener('change', function () {
        const checkboxes = document.querySelectorAll('.user-checkbox');
        checkboxes.forEach(checkbox => {
            checkbox.checked = selectAllCheckbox.checked;
        });
    });

    // 初始化加载用户数据
    loadUsers();
    
});