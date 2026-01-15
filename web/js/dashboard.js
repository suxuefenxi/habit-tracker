document.addEventListener('DOMContentLoaded', () => {
    Api.checkAuth();
    loadDashboard();
    loadUserStats();
});

async function loadDashboard() {
    try {
        // Fetch all habits. Ideally, there would be an endpoint for "today's habits"
        // For now, we list all active habits.
        const response = await Api.get('/habits');
        const habits = response.data || [];
        const activeHabits = habits.filter(h => h.is_active);

        const progressMap = await getHabitsProgress(activeHabits);
        renderHabits(activeHabits, progressMap);
    } catch (error) {
        console.error('Failed to load habits:', error);
    }
}

async function loadUserStats() {
    try {
        const response = await Api.get('/user/stats');
        const stats = response.data;
        
        if (stats) {
            document.getElementById('total-checkins').textContent = stats.total_checkins || 0;
            document.getElementById('current-streak').textContent = stats.longest_streak || 0; // Using longest streak as proxy if current not avail
            document.getElementById('total-points').textContent = stats.total_points || 0;
        }
    } catch (error) {
        console.error('Failed to load stats:', error);
    }
}

function renderHabits(habits, progressMap = {}) {
    const container = document.getElementById('habits-container');
    container.innerHTML = '';

    if (habits.length === 0) {
        container.innerHTML = '<div class="col-12"><p class="text-muted text-center">暂无进行中的习惯，去创建一个吧！</p></div>';
        return;
    }

    habits.forEach(habit => {
        const card = document.createElement('div');
        card.className = 'col-md-6 col-lg-4 mb-4';

        const progress = progressMap[habit.id] || { count: 0, target: habit.target_times || 1, percent: 0, reached: false };
        const percentText = `${progress.percent}%`;
        const btnDisabled = progress.reached;
        const btnClass = btnDisabled ? 'btn-secondary' : 'btn-primary';
        const btnText = btnDisabled ? '已完成' : '打卡 +1';
        const progressBarClass = progress.reached ? 'bg-success' : '';

        card.innerHTML = `
            <div class="card h-100 shadow-sm">
                <div class="card-body">
                    <h5 class="card-title">${habit.name}</h5>
                    <p class="card-text text-muted">${habit.description || '无描述'}</p>
                    <div class="mb-3">
                        <div class="d-flex justify-content-between align-items-center mb-2">
                            <small class="text-muted">今日进度</small>
                            <small class="fw-bold ${progress.reached ? 'text-success' : 'text-primary'}">${progress.count}/${progress.target} (${percentText})</small>
                        </div>
                        <div class="progress" style="height: 20px; background-color: #e9ecef;">
                            <div class="progress-bar ${progressBarClass}" role="progressbar" style="width: ${progress.percent}%" aria-valuenow="${progress.percent}" aria-valuemin="0" aria-valuemax="100">
                                ${progress.percent > 0 ? progress.percent + '%' : ''}
                            </div>
                        </div>
                    </div>
                    <div class="d-flex justify-content-between align-items-center mt-3">
                        <span class="badge bg-info text-dark">${habit.target_type === 'daily' ? '每天' : '每周'} ${habit.target_times} 次</span>
                        <button class="btn ${btnClass} btn-sm checkin-btn" data-id="${habit.id}" ${btnDisabled ? 'disabled' : ''}>
                            ${btnText}
                        </button>
                    </div>
                </div>
            </div>
        `;
        container.appendChild(card);
    });

    // Add event listeners to buttons
    document.querySelectorAll('.checkin-btn').forEach(btn => {
        btn.addEventListener('click', handleCheckin);
    });
}

async function handleCheckin(e) {
    const btn = e.target;
    const habitId = btn.dataset.id;
    
    // Disable button to prevent double clicks
    btn.disabled = true;
    
    try {
        await Api.post('/checkins', {
            habit_id: parseInt(habitId),
            count: 1
        });
        
        // Show success feedback (maybe a toast or just alert)
        // alert('打卡成功！');
        
        // Refresh stats and dashboard (progress + button state)
        await loadUserStats();
        await loadDashboard();
        
    } catch (error) {
        alert('打卡失败: ' + error.message);
        btn.disabled = false;
    }
}

async function getHabitsProgress(habits) {
    const results = await Promise.all(habits.map(habit => getHabitProgress(habit)));
    return results.reduce((acc, item) => {
        acc[item.id] = item.progress;
        return acc;
    }, {});
}

async function getHabitProgress(habit) {
    try {
        const response = await Api.get(`/habits/${habit.id}/checkins`);
        const checkins = response.data || [];
        const target = habit.target_times || 1;

        const now = new Date();
        const startOfToday = new Date(now.getFullYear(), now.getMonth(), now.getDate());
        const startOfWeek = getStartOfWeek(now);

        let count = 0;
        checkins.forEach(c => {
            const date = new Date(c.checkin_date);
            if (habit.target_type === 'daily') {
                if (date >= startOfToday && date < addDays(startOfToday, 1)) {
                    count += c.count || 0;
                }
            } else {
                if (date >= startOfWeek && date < addDays(startOfWeek, 7)) {
                    count += c.count || 0;
                }
            }
        });

        const percent = Math.min(100, Math.round((count / target) * 100));
        return {
            id: habit.id,
            progress: {
                count,
                target,
                percent,
                reached: count >= target
            }
        };
    } catch (error) {
        console.error('Failed to load habit progress:', error);
        return {
            id: habit.id,
            progress: {
                count: 0,
                target: habit.target_times || 1,
                percent: 0,
                reached: false
            }
        };
    }
}

function getStartOfWeek(date) {
    const d = new Date(date);
    const day = d.getDay();
    const diff = (day === 0 ? -6 : 1) - day; // Monday as start
    d.setDate(d.getDate() + diff);
    return new Date(d.getFullYear(), d.getMonth(), d.getDate());
}

function addDays(date, days) {
    const d = new Date(date);
    d.setDate(d.getDate() + days);
    return d;
}
