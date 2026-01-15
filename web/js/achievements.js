document.addEventListener('DOMContentLoaded', () => {
    Api.checkAuth();
    loadAchievements();
});

async function loadAchievements() {
    try {
        const allResponse = await Api.get('/achievements');
        const userResponse = await Api.get('/achievements/user');
        
        const allAchievements = allResponse.data || [];
        const userAchievements = userResponse.data || [];
        
        const unlockedMap = new Map();
        userAchievements.forEach(ua => {
            // Handle Go struct field casing (JSON usually lowercase but let's be safe)
            // Assuming JSON tags are lowercase snake_case or camelCase.
            // Go default JSON is struct field name if no tag, but usually we use tags.
            // Let's assume standard snake_case from GORM/JSON tags.
            const achId = ua.achievement_id || ua.AchievementID;
            const unlockedAt = ua.unlocked_at || ua.UnlockedAt;
            unlockedMap.set(achId, unlockedAt);
        });
        
        renderAchievements(allAchievements, unlockedMap);

    } catch (error) {
        console.error('Failed to load achievements:', error);
        const container = document.getElementById('achievements-container');
        if (container) {
            container.innerHTML = '<p class="text-danger">加载失败，请稍后再试</p>';
        }
    }
}

function renderAchievements(all, unlockedMap) {
    const container = document.getElementById('achievements-container');
    if (!container) return;
    
    container.innerHTML = '';
    
    if (all.length === 0) {
        container.innerHTML = '<p>暂无成就配置</p>';
        return;
    }

    all.forEach(ach => {
        const id = ach.id || ach.ID;
        const isUnlocked = unlockedMap.has(id);
        const unlockedAt = unlockedMap.get(id);
        
        const card = document.createElement('div');
        card.className = 'col-md-6 col-lg-4 mb-4';
        
        card.innerHTML = `
            <div class="card h-100 ${isUnlocked ? 'border-success' : 'border-secondary'}">
                <div class="card-body text-center">
                    <div class="display-4 mb-3">
                        ${isUnlocked ? '🏆' : '🔒'}
                    </div>
                    <h5 class="card-title">${ach.name || ach.Name}</h5>
                    <p class="card-text text-muted">${ach.description || ach.Description}</p>
                    ${isUnlocked ? 
                        `<span class="badge bg-success">已解锁: ${new Date(unlockedAt).toLocaleDateString()}</span>` : 
                        `<span class="badge bg-secondary">未解锁</span>`
                    }
                </div>
            </div>
        `;
        container.appendChild(card);
    });
}
