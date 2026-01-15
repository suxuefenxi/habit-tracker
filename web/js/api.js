const API_BASE_URL = '/api/v1';

class Api {
    static getToken() {
        return localStorage.getItem('token');
    }

    static setToken(token) {
        localStorage.setItem('token', token);
    }

    static removeToken() {
        localStorage.removeItem('token');
    }

    static async request(endpoint, method = 'GET', body = null) {
        const headers = {
            'Content-Type': 'application/json',
        };

        const token = this.getToken();
        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }

        const config = {
            method,
            headers,
        };

        if (body) {
            config.body = JSON.stringify(body);
        }

        try {
            const response = await fetch(`${API_BASE_URL}${endpoint}`, config);
            
            if (response.status === 401) {
                // Token expired or invalid
                // Only redirect if we are not already on the login or register page
                if (!window.location.pathname.includes('login.html') && !window.location.pathname.includes('register.html')) {
                    this.removeToken();
                    window.location.href = '/static/login.html';
                }
                throw new Error('Unauthorized');
            }

            const data = await response.json();
            
            if (!response.ok) {
                throw new Error(data.message || 'Something went wrong');
            }

            return data;
        } catch (error) {
            console.error('API Error:', error);
            throw error;
        }
    }

    static get(endpoint) {
        return this.request(endpoint, 'GET');
    }

    static post(endpoint, body) {
        return this.request(endpoint, 'POST', body);
    }

    static put(endpoint, body) {
        return this.request(endpoint, 'PUT', body);
    }

    static patch(endpoint, body) {
        return this.request(endpoint, 'PATCH', body);
    }

    static delete(endpoint) {
        return this.request(endpoint, 'DELETE');
    }
    
    static checkAuth() {
        if (!this.getToken()) {
            window.location.href = '/static/login.html';
        }
    }
}
