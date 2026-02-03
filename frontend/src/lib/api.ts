const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

interface ApiResponse<T = unknown> {
  code?: number
  msg: string
  data?: T
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

  if (!response.ok) {
    if (response.status === 401) {
      throw new Error('未登录或登录已过期')
    }
    const data: ApiResponse<T> = await response.json()
    throw new Error(data.msg || '请求失败')
  }

  const data: ApiResponse<T> = await response.json()

  if (data.code && data.code !== 0) {
    throw new Error(data.msg || '请求失败')
  }

  return data.data as T
}

export async function login(email: string, password: string): Promise<void> {
  const response = await fetch(`${API_BASE_URL}/users/login`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include',
    body: JSON.stringify({ email, password }),
  })

  const token = response.headers.get('x-jwt-token')
  if (token) {
    localStorage.setItem('access_token', token)
  }

  if (!response.ok) {
    throw new Error('登录失败')
  }
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

export interface UserProfile {
  id: number
  email: string
  nickname: string
  birthday: number
  about_me: string
  phone: string
}

export async function getUserProfile(): Promise<UserProfile> {
  return request<UserProfile>('/users/profile', {
    method: 'GET',
  })
}

export async function updateUserProfile(data: Partial<UserProfile>): Promise<void> {
  await request<void>('/users/profile', {
    method: 'PUT',
    body: JSON.stringify(data),
  })
}
