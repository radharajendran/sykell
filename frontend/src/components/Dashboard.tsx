import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { PlusIcon, PlayIcon, StopIcon, TrashIcon } from '@heroicons/react/24/outline'
import { MagnifyingGlassIcon } from '@heroicons/react/24/solid'
import UrlForm from './UrlForm'
import UrlTable from './UrlTable'
import type { CrawlJob } from '../types'
import { crawlerAPI, authAPI } from '../services/api'
import { transformCrawlURL } from '../utils/dataTransform'

const Dashboard = () => {
    const navigate = useNavigate()
    const [urls, setUrls] = useState<CrawlJob[]>([])
    const [showAddForm, setShowAddForm] = useState(false)
    const [searchTerm, setSearchTerm] = useState('')
    const [selectedUrls, setSelectedUrls] = useState<string[]>([])
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        fetchUrls()
    }, [])

    // Refresh URLs when search term changes (with debounce)
    useEffect(() => {
        const timeoutId = setTimeout(() => {
            if (!loading) {
                fetchUrls()
            }
        }, 500)

        return () => clearTimeout(timeoutId)
    }, [searchTerm])

    const fetchUrls = async () => {
        try {
            setLoading(true)

            // Check if user is authenticated
            if (!authAPI.isAuthenticated()) {
                navigate('/login')
                return
            }

            const response = await crawlerAPI.getCrawlURLs({
                page: 1,
                limit: 100, // Adjust as needed
                search: searchTerm
            })

            // Transform backend response to frontend format
            const transformedUrls: CrawlJob[] = response.data.map(transformCrawlURL)

            setUrls(transformedUrls)
        } catch (error) {
            console.error('Failed to fetch URLs:', error)
            // If authentication failed, redirect to login
            if (error instanceof Error && error.message === 'Authentication required') {
                navigate('/login')
            }
        } finally {
            setLoading(false)
        }
    }

    const handleAddUrl = async (url: string) => {
        try {
            // Add URL to backend
            await crawlerAPI.addURL({ url })

            // Refresh the URLs list
            await fetchUrls()

            setShowAddForm(false)

            console.log('URL added successfully:', url)
        } catch (error) {
            console.error('Failed to add URL:', error)
            // You could show an error toast here
        }
    }

    const handleBulkAction = async (action: 'start' | 'stop' | 'delete') => {
        if (selectedUrls.length === 0) return

        try {
            switch (action) {
                case 'start':
                    // Start crawling selected URLs
                    await Promise.all(
                        selectedUrls.map(id => crawlerAPI.startCrawl(id))
                    )
                    console.log('Started crawl jobs:', selectedUrls)
                    break
                case 'stop':
                    // TODO: Implement stop functionality in backend
                    console.log('Stopping crawl jobs:', selectedUrls)
                    break
                case 'delete':
                    await crawlerAPI.deleteURLs(selectedUrls)
                    console.log('Deleted crawl jobs:', selectedUrls)
                    break
            }

            // Refresh the URLs list and clear selection
            await fetchUrls()
            setSelectedUrls([])
        } catch (error) {
            console.error(`Failed to ${action} URLs:`, error)
            // You could show an error toast here
        }
    }

    const filteredUrls = urls.filter(url =>
        url.url.toLowerCase().includes(searchTerm.toLowerCase()) ||
        url.title.toLowerCase().includes(searchTerm.toLowerCase())
    )

    return (
        <div className="min-h-screen bg-gray-50">
            {/* Header */}
            <div className="bg-white border-b border-gray-200">
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                    <div className="flex justify-between items-center py-6">
                        <div>
                            <h1 className="text-3xl font-bold text-gray-900">Web Crawler Dashboard</h1>
                            <p className="mt-1 text-sm text-gray-500">
                                Analyze websites and track crawling progress
                            </p>
                        </div>
                        <div className="flex space-x-3">
                            <button
                                onClick={() => setShowAddForm(true)}
                                className="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                            >
                                <PlusIcon className="w-5 h-5 mr-2" />
                                Add URL
                            </button>
                            <button
                                onClick={() => {
                                    authAPI.logout()
                                    navigate('/login')
                                }}
                                className="inline-flex items-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                            >
                                Logout
                            </button>
                        </div>
                    </div>
                </div>
            </div>

            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
                {/* Search and Filters */}
                <div className="mb-6">
                    <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
                        <div className="relative">
                            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                                <MagnifyingGlassIcon className="h-5 w-5 text-gray-400" />
                            </div>
                            <input
                                type="text"
                                placeholder="Search URLs..."
                                value={searchTerm}
                                onChange={(e) => setSearchTerm(e.target.value)}
                                className="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md leading-5 bg-white placeholder-gray-500 focus:outline-none focus:placeholder-gray-400 focus:ring-1 focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                            />
                        </div>

                        {/* Bulk Actions */}
                        {selectedUrls.length > 0 && (
                            <div className="flex space-x-2">
                                <button
                                    onClick={() => handleBulkAction('start')}
                                    className="inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                                >
                                    <PlayIcon className="w-4 h-4 mr-1" />
                                    Start ({selectedUrls.length})
                                </button>
                                <button
                                    onClick={() => handleBulkAction('stop')}
                                    className="inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                                >
                                    <StopIcon className="w-4 h-4 mr-1" />
                                    Stop
                                </button>
                                <button
                                    onClick={() => handleBulkAction('delete')}
                                    className="inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm leading-4 font-medium rounded-md text-red-700 bg-white hover:bg-red-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
                                >
                                    <TrashIcon className="w-4 h-4 mr-1" />
                                    Delete
                                </button>
                            </div>
                        )}
                    </div>
                </div>

                {/* URL Table */}
                <UrlTable
                    urls={filteredUrls}
                    selectedUrls={selectedUrls}
                    onSelectionChange={setSelectedUrls}
                    onViewDetails={(id: string) => navigate(`/url/${id}`)}
                    loading={loading}
                />
            </div>

            {/* Add URL Modal */}
            {showAddForm && (
                <UrlForm
                    onSubmit={handleAddUrl}
                    onCancel={() => setShowAddForm(false)}
                />
            )}
        </div>
    )
}

export default Dashboard
