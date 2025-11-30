import { Link, useNavigate } from 'react-router-dom'
import { useEffect, useState } from 'react'

export default function Navbar() {
  const [isAuthenticated, setIsAuthenticated] = useState(false)
  const navigate = useNavigate()

  useEffect(() => {
    const token = localStorage.getItem('token')
    setIsAuthenticated(!!token)
  }, [])

  const handleLogout = () => {
    localStorage.removeItem('token')
    setIsAuthenticated(false)
    navigate('/home')
  }

  return (
    <nav className="bg-blue-600 text-white shadow-lg">
      <div className="container mx-auto px-4 py-3">
        <div className="flex justify-between items-center">
          <div className="flex space-x-4">
            <Link to="/home" className="hover:text-blue-200">Home</Link>
            {isAuthenticated && (
              <>
                <Link to="/projects" className="hover:text-blue-200">Projects</Link>
                <Link to="/sandbox" className="hover:text-blue-200">Sandbox</Link>
              </>
            )}
          </div>
          <div className="flex space-x-4">
            {!isAuthenticated ? (
              <>
                <Link to="/registration" className="hover:text-blue-200">Register</Link>
                <Link to="/authorization" className="hover:text-blue-200">Login</Link>
              </>
            ) : (
              <button onClick={handleLogout} className="hover:text-blue-200">Logout</button>
            )}
          </div>
        </div>
      </div>
    </nav>
  )
}

