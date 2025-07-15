const API_BASE_URL = 'http://localhost:8080/api'

export interface LoginRequest {
    email: string
    password: string
}

export interface LoginResponse {
    token: string
    user: {
        id: number
        name: string
        email: string
    }
}

export interface CreateUserRequest {
    name: string
    email: string
    password: string
}

export interface AddURLRequest {
    url: string
}

export interface BulkAddURLsRequest {
    urls: string[]
}

export interface CrawlURLsParams {
    page?: number
    limit?: number
    status?: string
    search?: string
}

// Store JWT token
let authToken: string | null = localStorage.getItem('authToken')

// Helper function to make authenticated requests
const makeRequest = async (url: string, options: RequestInit = {}) => {
    const headers: Record<string, string> = {
        'Content-Type': 'application/json',
        ...(options.headers as Record<string, string>),
    }

    if (authToken) {
        headers['Authorization'] = `Bearer ${authToken}`
    }

    const response = await fetch(`${API_BASE_URL}${url}`, {
        ...options,
        headers,
    })

    if (response.status === 401) {
        // Token expired or invalid
        authToken = null
        localStorage.removeItem('authToken')
        throw new Error('Authentication required')
    }

    if (!response.ok) {
        const errorData = await response.json().catch(() => ({}))
        throw new Error(errorData.error || 'Request failed')
    }

    return response.json()
}

// Authentication API
export const authAPI = {
    login: async (credentials: LoginRequest): Promise<LoginResponse> => {
        const response = await fetch(`${API_BASE_URL}/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(credentials),
        })

        if (!response.ok) {
            const errorData = await response.json().catch(() => ({}))
            throw new Error(errorData.error || 'Login failed')
        }

        const data = await response.json()
        authToken = data.token
        localStorage.setItem('authToken', data.token)
        return data
    },

    register: async (userData: CreateUserRequest): Promise<LoginResponse> => {
        const response = await fetch(`${API_BASE_URL}/users`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(userData),
        })

        if (!response.ok) {
            const errorData = await response.json().catch(() => ({}))
            throw new Error(errorData.error || 'Registration failed')
        }

        return response.json()
    },

    logout: () => {
        authToken = null
        localStorage.removeItem('authToken')
    },

    isAuthenticated: () => !!authToken,
}

// Crawler API
export const crawlerAPI = {
    // Add single URL
    addURL: async (request: AddURLRequest) => {
        return makeRequest('/crawler/urls', {
            method: 'POST',
            body: JSON.stringify(request),
        })
    },

    // Add multiple URLs
    bulkAddURLs: async (request: BulkAddURLsRequest) => {
        return makeRequest('/crawler/urls/bulk', {
            method: 'POST',
            body: JSON.stringify(request),
        })
    },

    // Get crawl URLs with pagination and filtering
    getCrawlURLs: async (params: CrawlURLsParams = {}) => {
        const queryParams = new URLSearchParams()

        if (params.page) queryParams.append('page', params.page.toString())
        if (params.limit) queryParams.append('limit', params.limit.toString())
        if (params.status) queryParams.append('status', params.status)
        if (params.search) queryParams.append('search', params.search)

        const url = `/crawler/urls${queryParams.toString() ? '?' + queryParams.toString() : ''}`
        return makeRequest(url)
    },

    // Get detailed crawl result
    getCrawlResult: async (id: string) => {
        return makeRequest(`/crawler/urls/${id}`)
    },

    // Start crawling a URL
    startCrawl: async (id: string) => {
        return makeRequest(`/crawler/urls/${id}/crawl`, {
            method: 'POST',
        })
    },

    // Delete multiple URLs
    deleteURLs: async (ids: string[]) => {
        return makeRequest('/crawler/urls', {
            method: 'DELETE',
            body: JSON.stringify({ ids }),
        })
    },

    // Re-crawl multiple URLs
    reCrawlURLs: async (ids: string[]) => {
        return makeRequest('/crawler/urls/recrawl', {
            method: 'POST',
            body: JSON.stringify({ ids }),
        })
    },

    // Get crawl statistics
    getStats: async () => {
        return makeRequest('/crawler/stats')
    },
}

// Export the auth token setter for external use
export const setAuthToken = (token: string) => {
    authToken = token
    localStorage.setItem('authToken', token)
}

export const getAuthToken = () => authToken
