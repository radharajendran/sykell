export interface CrawlJob {
    id: string
    url: string
    title: string
    status: 'queued' | 'running' | 'completed' | 'error'
    htmlVersion: string
    internalLinks: number
    externalLinks: number
    brokenLinks: number
    hasLoginForm: boolean
    headingCounts: {
        h1: number
        h2: number
        h3: number
        h4: number
        h5: number
        h6: number
    }
    brokenLinkDetails?: {
        id: number
        url: string
        statusCode: number
        error?: string
    }[]
    createdAt: Date
    completedAt: Date | null
    duration?: number
    error?: string
}

// Backend response types
export interface BackendCrawlURL {
    id: number
    url: string
    status: string
    title: string
    html_version: string
    h1_count: number
    h2_count: number
    h3_count: number
    h4_count: number
    h5_count: number
    h6_count: number
    internal_links_count: number
    external_links_count: number
    inaccessible_links_count: number
    has_login_form: boolean
    error_message: string
    last_crawled_at: string | null
    created_at: string
    updated_at: string
}

export interface BackendBrokenLink {
    id: number
    crawl_url_id: number
    url: string
    status_code: number
    error_message: string
    created_at: string
}

export interface BackendCrawlResult {
    crawl_url: BackendCrawlURL
    broken_links: BackendBrokenLink[]
}

export interface ApiResponse<T> {
    data: T
    message?: string
    error?: string
}

export interface PaginationData {
    page: number
    limit: number
    total: number
    totalPages: number
}

export interface CrawlJobsResponse extends ApiResponse<CrawlJob[]> {
    pagination: PaginationData
}

export interface SocketEvent {
    type: 'crawl_progress' | 'crawl_completed' | 'crawl_failed'
    jobId: string
    data: any
}
