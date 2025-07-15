import { useState } from 'react'
import {
    CheckBadgeIcon,
    ClockIcon,
    ExclamationTriangleIcon,
    StopIcon,
    PlayIcon
} from '@heroicons/react/24/outline'
import type { CrawlJob } from '../types'

interface UrlTableProps {
    urls: CrawlJob[]
    selectedUrls: string[]
    onSelectionChange: (selected: string[]) => void
    onViewDetails?: (id: string) => void
    loading: boolean
}

const UrlTable = ({ urls, selectedUrls, onSelectionChange, onViewDetails, loading }: UrlTableProps) => {
    const [sortField, setSortField] = useState<keyof CrawlJob>('createdAt')
    const [sortDirection, setSortDirection] = useState<'asc' | 'desc'>('desc')

    const getStatusIcon = (status: CrawlJob['status']) => {
        switch (status) {
            case 'completed':
                return <CheckBadgeIcon className="w-5 h-5 text-green-500" />
            case 'running':
                return <ClockIcon className="w-5 h-5 text-blue-500 animate-spin" />
            case 'failed':
                return <ExclamationTriangleIcon className="w-5 h-5 text-red-500" />
            case 'stopped':
                return <StopIcon className="w-5 h-5 text-gray-500" />
            default:
                return <ClockIcon className="w-5 h-5 text-yellow-500" />
        }
    }

    const getStatusBadge = (status: CrawlJob['status']) => {
        const baseClasses = "inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"

        switch (status) {
            case 'completed':
                return `${baseClasses} bg-green-100 text-green-800`
            case 'running':
                return `${baseClasses} bg-blue-100 text-blue-800`
            case 'failed':
                return `${baseClasses} bg-red-100 text-red-800`
            case 'stopped':
                return `${baseClasses} bg-gray-100 text-gray-800`
            default:
                return `${baseClasses} bg-yellow-100 text-yellow-800`
        }
    }

    const handleSort = (field: keyof CrawlJob) => {
        if (sortField === field) {
            setSortDirection(sortDirection === 'asc' ? 'desc' : 'asc')
        } else {
            setSortField(field)
            setSortDirection('asc')
        }
    }

    const sortedUrls = [...urls].sort((a, b) => {
        const aValue = a[sortField]
        const bValue = b[sortField]

        // Handle null/undefined values
        if (aValue == null && bValue == null) return 0
        if (aValue == null) return sortDirection === 'asc' ? -1 : 1
        if (bValue == null) return sortDirection === 'asc' ? 1 : -1

        if (aValue < bValue) return sortDirection === 'asc' ? -1 : 1
        if (aValue > bValue) return sortDirection === 'asc' ? 1 : -1
        return 0
    })

    const handleSelectAll = (checked: boolean) => {
        if (checked) {
            onSelectionChange(urls.map(url => url.id))
        } else {
            onSelectionChange([])
        }
    }

    const handleSelectUrl = (urlId: string, checked: boolean) => {
        if (checked) {
            onSelectionChange([...selectedUrls, urlId])
        } else {
            onSelectionChange(selectedUrls.filter(id => id !== urlId))
        }
    }

    const formatDate = (date: Date) => {
        return new Intl.DateTimeFormat('en-US', {
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        }).format(date)
    }

    if (loading) {
        return (
            <div className="bg-white shadow overflow-hidden sm:rounded-md">
                <div className="px-4 py-5 sm:p-6">
                    <div className="animate-pulse">
                        <div className="h-4 bg-gray-200 rounded w-1/4 mb-4"></div>
                        <div className="space-y-3">
                            {[...Array(5)].map((_, i) => (
                                <div key={i} className="h-4 bg-gray-200 rounded"></div>
                            ))}
                        </div>
                    </div>
                </div>
            </div>
        )
    }

    if (urls.length === 0) {
        return (
            <div className="bg-white shadow overflow-hidden sm:rounded-md">
                <div className="px-4 py-5 sm:p-6 text-center">
                    <PlayIcon className="mx-auto h-12 w-12 text-gray-400" />
                    <h3 className="mt-2 text-sm font-medium text-gray-900">No URLs added</h3>
                    <p className="mt-1 text-sm text-gray-500">
                        Get started by adding a website URL to crawl.
                    </p>
                </div>
            </div>
        )
    }

    return (
        <div className="bg-white shadow overflow-hidden sm:rounded-md">
            <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                    <thead className="bg-gray-50">
                        <tr>
                            <th className="px-6 py-3 text-left">
                                <input
                                    type="checkbox"
                                    checked={selectedUrls.length === urls.length && urls.length > 0}
                                    onChange={(e) => handleSelectAll(e.target.checked)}
                                    className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                                />
                            </th>
                            <th
                                className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
                                onClick={() => handleSort('url')}
                            >
                                URL
                            </th>
                            <th
                                className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
                                onClick={() => handleSort('title')}
                            >
                                Title
                            </th>
                            <th
                                className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
                                onClick={() => handleSort('status')}
                            >
                                Status
                            </th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                HTML Version
                            </th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                Links
                            </th>
                            <th
                                className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
                                onClick={() => handleSort('createdAt')}
                            >
                                Created
                            </th>
                        </tr>
                    </thead>
                    <tbody className="bg-white divide-y divide-gray-200">
                        {sortedUrls.map((url) => (
                            <tr
                                key={url.id}
                                className="hover:bg-gray-50 cursor-pointer"
                                onClick={() => onViewDetails?.(url.id)}
                            >
                                <td className="px-6 py-4 whitespace-nowrap">
                                    <input
                                        type="checkbox"
                                        checked={selectedUrls.includes(url.id)}
                                        onChange={(e) => {
                                            e.stopPropagation()
                                            handleSelectUrl(url.id, e.target.checked)
                                        }}
                                        className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                                    />
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap">
                                    <div className="flex items-center">
                                        {getStatusIcon(url.status)}
                                        <div className="ml-3">
                                            <div className="text-sm font-medium text-gray-900 truncate max-w-xs">
                                                {url.url}
                                            </div>
                                        </div>
                                    </div>
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap">
                                    <div className="text-sm text-gray-900 truncate max-w-xs">
                                        {url.title || '-'}
                                    </div>
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap">
                                    <span className={getStatusBadge(url.status)}>
                                        {url.status}
                                    </span>
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                    {url.htmlVersion || '-'}
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                    <div className="space-y-1">
                                        <div>Internal: {url.internalLinks}</div>
                                        <div>External: {url.externalLinks}</div>
                                        {url.brokenLinks > 0 && (
                                            <div className="text-red-600">Broken: {url.brokenLinks}</div>
                                        )}
                                    </div>
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                    {formatDate(url.createdAt)}
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    )
}

export default UrlTable
