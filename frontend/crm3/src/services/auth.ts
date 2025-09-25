// 認証API サービス
export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  email: string
  password: string
  name: string
}

export interface PasswordResetRequest {
  email: string
}

export interface PasswordResetConfirmRequest {
  token: string
  password: string
}

export interface AuthResponse {
  token: string
  user: {
    id: number
    email: string
    name: string
  }
}

export interface UserResponse {
  id: number
  email: string
  name: string
}

export interface MessageResponse {
  message: string
}

export interface ErrorResponse {
  error: string
  message?: string
}

const API_BASE_URL = 'http://localhost:1323/api'

class AuthService {
  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${API_BASE_URL}${endpoint}`

    const config: RequestInit = {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      ...options,
    }

    // 認証トークンがあれば追加
    const token = localStorage.getItem('auth_token')
    if (token) {
      config.headers = {
        ...config.headers,
        Authorization: `Bearer ${token}`,
      }
    }

    const response = await fetch(url, config)

    if (!response.ok) {
      const errorData: ErrorResponse = await response.json()
      throw new Error(errorData.message || errorData.error || `HTTP ${response.status}`)
    }

    return response.json()
  }

  // ユーザー登録
  async register(data: RegisterRequest): Promise<AuthResponse> {
    return this.request<AuthResponse>('/auth/register', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  // ログイン
  async login(data: LoginRequest): Promise<AuthResponse> {
    return this.request<AuthResponse>('/auth/login', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  // パスワードリセット要求
  async requestPasswordReset(data: PasswordResetRequest): Promise<MessageResponse> {
    return this.request<MessageResponse>('/auth/reset-password', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  // パスワードリセット確認
  async confirmPasswordReset(data: PasswordResetConfirmRequest): Promise<MessageResponse> {
    return this.request<MessageResponse>('/auth/confirm-reset', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  // ユーザープロフィール取得
  async getProfile(): Promise<UserResponse> {
    return this.request<UserResponse>('/user/profile')
  }

  // トークンをローカルストレージに保存
  saveToken(token: string): void {
    localStorage.setItem('auth_token', token)
  }

  // トークンをローカルストレージから取得
  getToken(): string | null {
    return localStorage.getItem('auth_token')
  }

  // トークンを削除（ログアウト）
  removeToken(): void {
    localStorage.removeItem('auth_token')
  }

  // 認証状態チェック
  isAuthenticated(): boolean {
    return !!this.getToken()
  }
}

export const authService = new AuthService()