const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

interface ApiResponse<T = unknown> {
  Code?: number
  Msg: string
  Data?: T
}

async function request<T>(url: string, options?: RequestInit): Promise<T> {
  const token = localStorage.getItem('access_token')
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  }

  if (options?.headers) {
    const h = options.headers as Record<string, string>
    Object.assign(headers, h)
  }

  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  const response = await fetch(`${API_BASE_URL}${url}`, {
    ...options,
    headers,
    credentials: 'include',
  })

  const data: ApiResponse<T> = await response.json()

  if (data.Code && data.Code !== 0) {
    throw new Error(data.Msg || '请求失败')
  }

  if (response.ok) {
    return data.Data as T
  }

  throw new Error(data.Msg || '请求失败')
}

export async function login(email: string, password: string): Promise<void> {
  await request<void>('/users/login', {
    method: 'POST',
    body: JSON.stringify({ email, password }),
  })
}

export async function register(email: string, password: string, confirmPassword: string, nickname: string): Promise<void> {
  await request<void>('/users/signup', {
    method: 'POST',
    body: JSON.stringify({
      email,
      password,
      confirm_password: confirmPassword,
      nickname,
    }),
  })
}

export async function logout(): Promise<void> {
  await request<void>('/users/logout', {
    method: 'POST',
  })
  localStorage.removeItem('access_token')
}

export async function refreshToken(): Promise<void> {
  await request<void>('/users/refresh_token', {
    method: 'GET',
  })
}

export function getAccessToken(): string | null {
  return localStorage.getItem('access_token')
}

export function setAccessToken(token: string): void {
  localStorage.setItem('access_token', token)
}

export function clearAccessToken(): void {
  localStorage.removeItem('access_token')
}
