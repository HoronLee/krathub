const API_BASE_URL = 'http://localhost:9000'

export const api = {
  async login(loginId: string, password: string) {
    const response = await fetch(`${API_BASE_URL}/v1/auth/login/using-email-password`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ loginId, password }),
    })
    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.message || `登录失败 (${response.status})`)
    }
    return response.json()
  },

  async signup(name: string, email: string, password: string, passwordConfirm: string) {
    const response = await fetch(`${API_BASE_URL}/v1/auth/signup/using-email`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name, email, password, passwordConfirm }),
    })
    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.message || `注册失败 (${response.status})`)
    }
    return response.json()
  },

  async getCurrentUser(token: string) {
    const response = await fetch(`${API_BASE_URL}/v1/user/info`, {
      method: 'GET',
      headers: { Authorization: `Bearer ${token}` },
    })
    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.message || `获取用户信息失败 (${response.status})`)
    }
    return response.json()
  },
}
