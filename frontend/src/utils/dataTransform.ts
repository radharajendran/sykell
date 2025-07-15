import type { CrawlJob, BackendCrawlURL, BackendCrawlResult } from '../types'

export const transformCrawlURL = (backendData: BackendCrawlURL): CrawlJob => {
    return {
        id: backendData.id.toString(),
        url: backendData.url,
        title: backendData.title || '',
        status: backendData.status as CrawlJob['status'],
        htmlVersion: backendData.html_version || '',
        internalLinks: backendData.internal_links_count || 0,
        externalLinks: backendData.external_links_count || 0,
        brokenLinks: backendData.inaccessible_links_count || 0,
        hasLoginForm: backendData.has_login_form || false,
        headingCounts: {
            h1: backendData.h1_count || 0,
            h2: backendData.h2_count || 0,
            h3: backendData.h3_count || 0,
            h4: backendData.h4_count || 0,
            h5: backendData.h5_count || 0,
            h6: backendData.h6_count || 0,
        },
        createdAt: new Date(backendData.created_at),
        completedAt: backendData.last_crawled_at ? new Date(backendData.last_crawled_at) : null,
        error: backendData.error_message || undefined,
        duration: backendData.last_crawled_at && backendData.created_at ?
            new Date(backendData.last_crawled_at).getTime() - new Date(backendData.created_at).getTime() :
            undefined
    }
}

export const transformCrawlResult = (backendData: BackendCrawlResult): CrawlJob => {
    const baseJob = transformCrawlURL(backendData.crawl_url)

    return {
        ...baseJob,
        brokenLinkDetails: backendData.broken_links.map(link => ({
            id: link.id,
            url: link.url,
            statusCode: link.status_code,
            error: link.error_message || undefined
        }))
    }
}
