document.addEventListener('DOMContentLoaded', () => {
    Api.checkAuth();
    loadLeaderboard('weekly');

    const tabWeekly = document.getElementById('tab-weekly');
    const tabMonthly = document.getElementById('tab-monthly');

    if (tabWeekly) {
        tabWeekly.addEventListener('click', (e) => {
            e.preventDefault();
            setActiveTab('weekly');
            loadLeaderboard('weekly');
        });
    }

    if (tabMonthly) {
        tabMonthly.addEventListener('click', (e) => {
            e.preventDefault();
            setActiveTab('monthly');
            loadLeaderboard('monthly');
        });
    }
});

function setActiveTab(type) {
    document.querySelectorAll('.nav-tabs .nav-link').forEach(el => el.classList.remove('active'));
    document.getElementById(`tab-${type}`).classList.add('active');
}

async function loadLeaderboard(type) {
    const tbody = document.getElementById('leaderboard-body');
    if (!tbody) return;
    
    tbody.innerHTML = '<tr><td colspan="3" class="text-center">加载中...</td></tr>';

    try {
        const response = await Api.get(`/leaderboard/${type}`);
        const data = response.data || [];
        
        tbody.innerHTML = '';
        
        if (data.length === 0) {
            tbody.innerHTML = '<tr><td colspan="3" class="text-center">暂无数据</td></tr>';
            return;
        }

        data.forEach((item, index) => {
            const tr = document.createElement('tr');
            let rankDisplay = index + 1;
            if (index === 0) rankDisplay = '🥇';
            else if (index === 1) rankDisplay = '🥈';
            else if (index === 2) rankDisplay = '🥉';
            
            tr.innerHTML = `
                <td class="fw-bold" style="font-size: 1.2rem;">${rankDisplay}</td>
                <td>${item.nickname || item.username || '用户' + item.user_id}</td>
                <td class="fw-bold text-primary">${item.points || 0}</td>
            `;
            tbody.appendChild(tr);
        });

    } catch (error) {
        console.error('Failed to load leaderboard:', error);
        tbody.innerHTML = '<tr><td colspan="3" class="text-center text-danger">加载失败</td></tr>';
    }
}
