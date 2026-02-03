import { createContext, useContext, useState, useEffect, type ReactNode } from 'react'
import { login, logout, register, getUserProfile, getAccessToken, clearAccessToken } from '@/lib/api'

export interface UserProfile {
  id: number
  email: string
  nickname: string
  birthday: number
  about_me: string
  phone: string
}

interface AuthContextType {
  isAuthenticated: boolean
  user: UserProfile | null
  isLoading: boolean
  login: (email: string, password: string) => Promise<void>
  register: (email: string, password: string, confirmPassword: string, nickname: string) => Promise<void>
  logout: () => Promise<void>
  refreshProfile: () => Promise<void>
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [isAuthenticated, setIsAuthenticated] = useState(false)
  const [user, setUser] = useState<UserProfile | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  const checkAuth = async () => {
    const token = getAccessToken()
    if (token) {
      try {
        const profile = await getUserProfile()
        setUser(profile)
        setIsAuthenticated(true)
      } catch (error) {
        clearAccessToken()
        setIsAuthenticated(false)
        setUser(null)
      }
    }
    setIsLoading(false)
  }

  useEffect(() => {
    checkAuth()
  }, [])

  const handleLogin = async (email: string, password: string) => {
    try {
      await login(email, password)
      const profile = await getUserProfile()
      setUser(profile)
      setIsAuthenticated(true)
    } catch (error) {
      clearAccessToken()
      setIsAuthenticated(false)
      setUser(null)
      throw error
    }
  }

  const handleRegister = async (email: string, password: string, confirmPassword: string, nickname: string) => {
    await register(email, password, confirmPassword, nickname)
  }

  const handleLogout = async () => {
    await logout()
    setIsAuthenticated(false)
    setUser(null)
  }

  const handleRefreshProfile = async () => {
    const profile = await getUserProfile()
    setUser(profile)
  }

  return (
    <AuthContext.Provider
      value={{
        isAuthenticated,
        user,
        isLoading,
        login: handleLogin,
        register: handleRegister,
        logout: handleLogout,
        refreshProfile: handleRefreshProfile,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
