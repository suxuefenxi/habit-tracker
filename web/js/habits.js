document.addEventListener('DOMContentLoaded', () => {
    Api.checkAuth();
    
    if (document.getElementById('habits-list')) {
        loadHabitsList();
        setupCreateHabit();
        setupEditHabit();
    }
    
    if (document.getElementById('habit-detail')) {
        loadHabitDetail();
    }
});

async function loadHabitsList() {
    try {
        const response = await Api.get('/habits');
        const habits = response.data || [];
        const tbody = document.getElementById('habits-table-body');
        tbody.innerHTML = '';
        
        habits.forEach(habit => {
            const tr = document.createElement('tr');
            tr.innerHTML = `
                <td>${habit.name}</td>
                <td>${habit.target_type === 'daily' ? '每天' : '每周'} ${habit.target_times} 次</td>
                <td><span class="badge ${habit.is_active ? 'bg-success' : 'bg-secondary'}">${habit.is_active ? '进行中' : '已停用'}</span></td>
                <td>
                    <a href="habit_detail.html?id=${habit.id}" class="btn btn-sm btn-info text-white">详情</a>
                    <button class="btn btn-sm btn-outline-secondary edit-btn" data-id="${habit.id}">编辑</button>
                </td>
            `;
            tbody.appendChild(tr);
        });

        // Add event listeners to edit buttons
        document.querySelectorAll('.edit-btn').forEach(btn => {
            btn.addEventListener('click', () => openEditModal(btn.dataset.id));
        });

    } catch (error) {
        console.error('Failed to load habits:', error);
    }
}

function setupCreateHabit() {
    const form = document.getElementById('createHabitForm');
    if (!form) return;

    form.addEventListener('submit', async (e) => {
        e.preventDefault();
        
        const data = {
            name: form.name.value,
            description: form.description.value,
            target_type: form.target_type.value,
            target_times: parseInt(form.target_times.value),
            start_date: new Date().toISOString().split('T')[0] // 只保留日期部分
        };

        try {
            await Api.post('/habits', data);
            // Close modal
            const modalEl = document.getElementById('createHabitModal');
            const modal = bootstrap.Modal.getInstance(modalEl);
            modal.hide();
            form.reset();
            loadHabitsList();
        } catch (error) {
            alert('创建失败: ' + error.message);
        }
    });
}

async function loadHabitDetail() {
    const urlParams = new URLSearchParams(window.location.search);
    const id = urlParams.get('id');
    
    if (!id) {
        alert('未指定习惯ID');
        window.location.href = 'habits.html';
        return;
    }

    try {
        const response = await Api.get(`/habits/${id}`);
        const habit = response.data;
        
        document.getElementById('habit-name').textContent = habit.name;
        document.getElementById('habit-desc').textContent = habit.description || '无描述';
        document.getElementById('habit-target').textContent = `${habit.target_type === 'daily' ? '每天' : '每周'} ${habit.target_times} 次`;
        document.getElementById('habit-status').innerHTML = `<span class="badge ${habit.is_active ? 'bg-success' : 'bg-secondary'}">${habit.is_active ? '进行中' : '已停用'}</span>`;

        // Load checkins
        loadHabitCheckins(id);

    } catch (error) {
        console.error('Failed to load habit detail:', error);
        alert('加载失败');
    }
}

async function loadHabitCheckins(id) {
    try {
        const response = await Api.get(`/habits/${id}/checkins`);
        const checkins = response.data || [];
        
        const tbody = document.getElementById('checkins-table-body');
        tbody.innerHTML = '';
        
        if (checkins.length === 0) {
            tbody.innerHTML = '<tr><td colspan="2" class="text-center">暂无打卡记录</td></tr>';
            return;
        }

        checkins.forEach(c => {
            const tr = document.createElement('tr');
            tr.innerHTML = `
                <td>${new Date(c.checkin_date).toLocaleDateString()}</td>
                <td>${c.count}</td>
            `;
            tbody.appendChild(tr);
        });
    } catch (error) {
        console.error('Failed to load checkins:', error);
    }
}

async function openEditModal(habitId) {
    try {
        const response = await Api.get(`/habits/${habitId}`);
        const habit = response.data;
        
        const form = document.getElementById('editHabitForm');
        form.id.value = habit.id;
        form.name.value = habit.name;
        form.description.value = habit.description || '';
        form.target_type.value = habit.target_type;
        form.target_times.value = habit.target_times;
        form.is_active.value = habit.is_active.toString();
        // Store start_date in YYYY-MM-DD format
        const startDate = habit.start_date.split('T')[0]; // Extract date part only
        form.dataset.startDate = startDate;
        
        const modal = new bootstrap.Modal(document.getElementById('editHabitModal'));
        modal.show();
    } catch (error) {
        alert('加载习惯信息失败: ' + error.message);
    }
}

function setupEditHabit() {
    const form = document.getElementById('editHabitForm');
    if (!form) return;

    form.addEventListener('submit', async (e) => {
        e.preventDefault();
        
        const habitId = form.id.value;
        const data = {
            name: form.name.value,
            description: form.description.value,
            target_type: form.target_type.value,
            target_times: parseInt(form.target_times.value),
            start_date: form.dataset.startDate, // Use stored start_date
            is_active: form.is_active.value === 'true'
        };

        try {
            await Api.put(`/habits/${habitId}`, data);
            const modalEl = document.getElementById('editHabitModal');
            const modal = bootstrap.Modal.getInstance(modalEl);
            modal.hide();
            form.reset();
            loadHabitsList();
        } catch (error) {
            alert('保存失败: ' + error.message);
        }
    });
}
