document.addEventListener('DOMContentLoaded', () => {
    const loginForm = document.getElementById('loginForm');
    const registerForm = document.getElementById('registerForm');

    if (loginForm) {
        loginForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const username = loginForm.username.value;
            const password = loginForm.password.value;

            try {
                const response = await Api.post('/auth/login', { username, password });
                // Assuming the response structure is { code: 200, message: "...", data: { token: "..." } }
                // Or if the backend returns the token directly in the data field or root.
                // Based on common practices and the task description, let's assume data contains the token.
                // If the backend returns { token: "..." } directly, adjust accordingly.
                // Let's assume the backend returns { token: "..." } or { data: { token: "..." } }
                
                // Checking the backend code would be ideal, but let's assume a standard response wrapper if used, 
                // or direct return. The task description says "Unified response structure: { code, message, data }".
                
                const token = response.data ? response.data.token : response.token;
                
                if (token) {
                    Api.setToken(token);
                    window.location.href = '/static/dashboard.html';
                } else {
                    throw new Error('Token not found in response');
                }
            } catch (error) {
                alert('Login failed: ' + error.message);
            }
        });
    }

    if (registerForm) {
        registerForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const username = registerForm.username.value;
            const password = registerForm.password.value;
            const nickname = registerForm.nickname.value;

            try {
                await Api.post('/auth/register', { username, password, nickname });
                alert('Registration successful! Please login.');
                window.location.href = '/static/login.html';
            } catch (error) {
                alert('Registration failed: ' + error.message);
            }
        });
    }
});
