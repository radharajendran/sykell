import { useState } from 'react'
import { XMarkIcon } from '@heroicons/react/24/outline'

interface UrlFormProps {
    onSubmit: (url: string) => void
    onCancel: () => void
}

const UrlForm = ({ onSubmit, onCancel }: UrlFormProps) => {
    const [url, setUrl] = useState('')
    const [isValid, setIsValid] = useState(true)
    const [error, setError] = useState('')

    const validateUrl = (input: string) => {
        try {
            const urlObj = new URL(input)
            return urlObj.protocol === 'http:' || urlObj.protocol === 'https:'
        } catch {
            return false
        }
    }

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault()

        if (!url.trim()) {
            setError('URL is required')
            setIsValid(false)
            return
        }

        if (!validateUrl(url)) {
            setError('Please enter a valid HTTP or HTTPS URL')
            setIsValid(false)
            return
        }

        onSubmit(url.trim())
    }

    const handleUrlChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const value = e.target.value
        setUrl(value)

        if (value && !validateUrl(value)) {
            setError('Please enter a valid HTTP or HTTPS URL')
            setIsValid(false)
        } else {
            setError('')
            setIsValid(true)
        }
    }

    return (
        <div className="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
            <div className="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white">
                <div className="flex justify-between items-center mb-4">
                    <h3 className="text-lg font-medium text-gray-900">Add New URL</h3>
                    <button
                        onClick={onCancel}
                        className="text-gray-400 hover:text-gray-600"
                    >
                        <XMarkIcon className="w-6 h-6" />
                    </button>
                </div>

                <form onSubmit={handleSubmit}>
                    <div className="mb-4">
                        <label htmlFor="url" className="block text-sm font-medium text-gray-700 mb-2">
                            Website URL
                        </label>
                        <input
                            type="url"
                            id="url"
                            value={url}
                            onChange={handleUrlChange}
                            placeholder="https://example.com"
                            className={`w-full px-3 py-2 border rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500 ${!isValid ? 'border-red-300' : 'border-gray-300'
                                }`}
                            autoFocus
                        />
                        {error && (
                            <p className="mt-1 text-sm text-red-600">{error}</p>
                        )}
                    </div>

                    <div className="flex justify-end space-x-3">
                        <button
                            type="button"
                            onClick={onCancel}
                            className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 border border-gray-300 rounded-md hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500"
                        >
                            Cancel
                        </button>
                        <button
                            type="submit"
                            disabled={!url.trim() || !isValid}
                            className="px-4 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
                        >
                            Add URL
                        </button>
                    </div>
                </form>
            </div>
        </div>
    )
}

export default UrlForm
