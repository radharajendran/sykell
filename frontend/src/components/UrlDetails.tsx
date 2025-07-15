import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { ArrowLeftIcon, LinkIcon, ExclamationTriangleIcon } from '@heroicons/react/24/outline'
import { PieChart, Pie, Cell, BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts'
import type { CrawlJob } from '../types'
import { crawlerAPI, authAPI } from '../services/api'
import { transformCrawlResult } from '../utils/dataTransform'

const UrlDetails = () => {
    const { id } = useParams<{ id: string }>()
    const navigate = useNavigate()
    const [crawlJob, setCrawlJob] = useState<CrawlJob | null>(null)
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        if (id) {
            fetchCrawlJob(id)
        }
    }, [id])

    const fetchCrawlJob = async (jobId: string) => {
        try {
            setLoading(true)

            // Check if user is authenticated
            if (!authAPI.isAuthenticated()) {
                navigate('/login')
                return
            }

            const response = await crawlerAPI.getCrawlResult(jobId)

            // Transform backend response to frontend format
            const transformedJob = transformCrawlResult(response)

            setCrawlJob(transformedJob)
        } catch (error) {
            console.error('Failed to fetch crawl job:', error)
            // If authentication failed, redirect to login
            if (error instanceof Error && error.message === 'Authentication required') {
                navigate('/login')
            }
        } finally {
            setLoading(false)
        }
    }

    if (loading) {
        return (
            <div className="min-h-screen bg-gray-50">
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
                    <div className="animate-pulse">
                        <div className="h-8 bg-gray-200 rounded w-1/4 mb-6"></div>
                        <div className="bg-white shadow rounded-lg p-6">
                            <div className="h-6 bg-gray-200 rounded w-1/3 mb-4"></div>
                            <div className="space-y-3">
                                {[...Array(8)].map((_, i) => (
                                    <div key={i} className="h-4 bg-gray-200 rounded"></div>
                                ))}
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        )
    }

    if (!crawlJob) {
        return (
            <div className="min-h-screen bg-gray-50">
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
                    <div className="text-center">
                        <h2 className="text-2xl font-bold text-gray-900">Crawl Job Not Found</h2>
                        <p className="mt-2 text-gray-600">The requested crawl job could not be found.</p>
                        <button
                            onClick={() => navigate('/')}
                            className="mt-4 inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700"
                        >
                            <ArrowLeftIcon className="w-4 h-4 mr-2" />
                            Back to Dashboard
                        </button>
                    </div>
                </div>
            </div>
        )
    }

    const linkData = [
        { name: 'Internal Links', value: crawlJob.internalLinks, color: '#3B82F6' },
        { name: 'External Links', value: crawlJob.externalLinks, color: '#10B981' },
        { name: 'Broken Links', value: crawlJob.brokenLinks, color: '#EF4444' },
    ]

    const headingData = Object.entries(crawlJob.headingCounts)
        .filter(([, count]) => count > 0)
        .map(([level, count]) => ({
            level: level.toUpperCase(),
            count
        }))

    const formatDuration = (duration: number) => {
        const seconds = Math.floor(duration / 1000)
        const minutes = Math.floor(seconds / 60)
        const remainingSeconds = seconds % 60

        if (minutes > 0) {
            return `${minutes}m ${remainingSeconds}s`
        }
        return `${remainingSeconds}s`
    }

    return (
        <div className="min-h-screen bg-gray-50">
            {/* Header */}
            <div className="bg-white border-b border-gray-200">
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                    <div className="flex items-center justify-between py-6">
                        <div className="flex items-center">
                            <button
                                onClick={() => navigate('/')}
                                className="mr-4 p-2 rounded-md text-gray-400 hover:text-gray-500 hover:bg-gray-100"
                            >
                                <ArrowLeftIcon className="w-6 h-6" />
                            </button>
                            <div>
                                <h1 className="text-2xl font-bold text-gray-900">Crawl Analysis</h1>
                                <p className="mt-1 text-sm text-gray-500">{crawlJob.url}</p>
                            </div>
                        </div>
                        <div className="flex items-center space-x-4">
                            <span className={`inline-flex items-center px-3 py-1 rounded-full text-sm font-medium ${crawlJob.status === 'completed' ? 'bg-green-100 text-green-800' :
                                crawlJob.status === 'running' ? 'bg-blue-100 text-blue-800' :
                                    crawlJob.status === 'error' ? 'bg-red-100 text-red-800' :
                                        'bg-yellow-100 text-yellow-800'
                                }`}>
                                {crawlJob.status}
                            </span>
                        </div>
                    </div>
                </div>
            </div>

            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
                <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                    {/* Overview Stats */}
                    <div className="lg:col-span-1">
                        <div className="bg-white shadow rounded-lg p-6">
                            <h3 className="text-lg font-medium text-gray-900 mb-4">Overview</h3>
                            <dl className="space-y-4">
                                <div>
                                    <dt className="text-sm font-medium text-gray-500">Title</dt>
                                    <dd className="mt-1 text-sm text-gray-900">{crawlJob.title || 'N/A'}</dd>
                                </div>
                                <div>
                                    <dt className="text-sm font-medium text-gray-500">HTML Version</dt>
                                    <dd className="mt-1 text-sm text-gray-900">{crawlJob.htmlVersion || 'N/A'}</dd>
                                </div>
                                <div>
                                    <dt className="text-sm font-medium text-gray-500">Has Login Form</dt>
                                    <dd className="mt-1 text-sm text-gray-900">
                                        {crawlJob.hasLoginForm ? 'Yes' : 'No'}
                                    </dd>
                                </div>
                                <div>
                                    <dt className="text-sm font-medium text-gray-500">Duration</dt>
                                    <dd className="mt-1 text-sm text-gray-900">
                                        {crawlJob.duration ? formatDuration(crawlJob.duration) : 'N/A'}
                                    </dd>
                                </div>
                                <div>
                                    <dt className="text-sm font-medium text-gray-500">Created</dt>
                                    <dd className="mt-1 text-sm text-gray-900">
                                        {crawlJob.createdAt.toLocaleDateString()} at {crawlJob.createdAt.toLocaleTimeString()}
                                    </dd>
                                </div>
                                {crawlJob.completedAt && (
                                    <div>
                                        <dt className="text-sm font-medium text-gray-500">Completed</dt>
                                        <dd className="mt-1 text-sm text-gray-900">
                                            {crawlJob.completedAt.toLocaleDateString()} at {crawlJob.completedAt.toLocaleTimeString()}
                                        </dd>
                                    </div>
                                )}
                            </dl>
                        </div>
                    </div>

                    {/* Charts Section */}
                    <div className="lg:col-span-2 space-y-6">
                        {/* Link Distribution Chart */}
                        <div className="bg-white shadow rounded-lg p-6">
                            <h3 className="text-lg font-medium text-gray-900 mb-4">Link Distribution</h3>
                            <div className="h-64">
                                <ResponsiveContainer width="100%" height="100%">
                                    <PieChart>
                                        <Pie
                                            data={linkData}
                                            cx="50%"
                                            cy="50%"
                                            labelLine={false}
                                            label={({ name, value }) => `${name}: ${value}`}
                                            outerRadius={80}
                                            fill="#8884d8"
                                            dataKey="value"
                                        >
                                            {linkData.map((entry, index) => (
                                                <Cell key={`cell-${index}`} fill={entry.color} />
                                            ))}
                                        </Pie>
                                        <Tooltip />
                                    </PieChart>
                                </ResponsiveContainer>
                            </div>
                        </div>

                        {/* Heading Distribution Chart */}
                        {headingData.length > 0 && (
                            <div className="bg-white shadow rounded-lg p-6">
                                <h3 className="text-lg font-medium text-gray-900 mb-4">Heading Distribution</h3>
                                <div className="h-64">
                                    <ResponsiveContainer width="100%" height="100%">
                                        <BarChart data={headingData}>
                                            <CartesianGrid strokeDasharray="3 3" />
                                            <XAxis dataKey="level" />
                                            <YAxis />
                                            <Tooltip />
                                            <Bar dataKey="count" fill="#3B82F6" />
                                        </BarChart>
                                    </ResponsiveContainer>
                                </div>
                            </div>
                        )}
                    </div>
                </div>

                {/* Broken Links Section */}
                {crawlJob.brokenLinks > 0 && crawlJob.brokenLinkDetails && (
                    <div className="mt-8">
                        <div className="bg-white shadow rounded-lg">
                            <div className="px-6 py-4 border-b border-gray-200">
                                <div className="flex items-center">
                                    <ExclamationTriangleIcon className="w-5 h-5 text-red-500 mr-2" />
                                    <h3 className="text-lg font-medium text-gray-900">
                                        Broken Links ({crawlJob.brokenLinks})
                                    </h3>
                                </div>
                            </div>
                            <div className="overflow-x-auto">
                                <table className="min-w-full divide-y divide-gray-200">
                                    <thead className="bg-gray-50">
                                        <tr>
                                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                                URL
                                            </th>
                                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                                Status Code
                                            </th>
                                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                                Error
                                            </th>
                                        </tr>
                                    </thead>
                                    <tbody className="bg-white divide-y divide-gray-200">
                                        {crawlJob.brokenLinkDetails.map((link, index) => (
                                            <tr key={index}>
                                                <td className="px-6 py-4 whitespace-nowrap">
                                                    <div className="flex items-center">
                                                        <LinkIcon className="w-4 h-4 text-gray-400 mr-2" />
                                                        <span className="text-sm text-gray-900 break-all">{link.url}</span>
                                                    </div>
                                                </td>
                                                <td className="px-6 py-4 whitespace-nowrap">
                                                    <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800">
                                                        {link.statusCode}
                                                    </span>
                                                </td>
                                                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                                    {link.error || '-'}
                                                </td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            </div>
                        </div>
                    </div>
                )}
            </div>
        </div>
    )
}

export default UrlDetails
