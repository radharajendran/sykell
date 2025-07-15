import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import Dashboard from './components/Dashboard'
import UrlDetails from './components/UrlDetails'
import Login from './components/Login'
import './App.css'

function App() {
  return (
    <Router>
      <div className="min-h-screen bg-gray-50">
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/" element={<Dashboard />} />
          <Route path="/url/:id" element={<UrlDetails />} />
        </Routes>
      </div>
    </Router>
  )
}

export default App
