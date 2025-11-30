import { useEffect, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import axios from 'axios'

interface Project {
  id: number
  name: string
  created_at: string
}

export default function Projects() {
  const [projects, setProjects] = useState<Project[]>([])
  const [newProjectName, setNewProjectName] = useState('')
  const [error, setError] = useState('')
  const navigate = useNavigate()

  useEffect(() => {
    loadProjects()
  }, [])

  const loadProjects = async () => {
    try {
      const token = localStorage.getItem('token')
      if (!token) {
        navigate('/authorization')
        return
      }

      const response = await axios.get('/api/projects', {
        headers: { Authorization: token },
      })
      setProjects(response.data)
    } catch (err: any) {
      if (err.response?.status === 401) {
        navigate('/authorization')
      }
    }
  }

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')

    try {
      const token = localStorage.getItem('token')
      const response = await axios.post(
        '/api/projects',
        { name: newProjectName },
        { headers: { Authorization: token } }
      )
      setProjects([...projects, response.data])
      setNewProjectName('')
      navigate(`/projects/${response.data.id}`)
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to create project')
    }
  }

  const handleDelete = async (id: number) => {
    if (!confirm('Are you sure you want to delete this project?')) return

    try {
      const token = localStorage.getItem('token')
      await axios.delete(`/api/projects/${id}`, {
        headers: { Authorization: token },
      })
      setProjects(projects.filter((p) => p.id !== id))
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to delete project')
    }
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="max-w-4xl mx-auto">
        <h1 className="text-3xl font-bold mb-6">Projects</h1>

        {error && <div className="bg-red-100 text-red-700 p-3 rounded mb-4">{error}</div>}

        <form onSubmit={handleCreate} className="mb-6">
          <div className="flex gap-2">
            <input
              type="text"
              value={newProjectName}
              onChange={(e) => setNewProjectName(e.target.value)}
              placeholder="Project name"
              className="flex-1 px-4 py-2 border rounded"
              required
            />
            <button type="submit" className="bg-blue-600 text-white px-6 py-2 rounded hover:bg-blue-700">
              Create
            </button>
          </div>
        </form>

        <div className="bg-white rounded-lg shadow">
          {projects.length === 0 ? (
            <div className="p-8 text-center text-gray-500">No projects yet. Create one to get started!</div>
          ) : (
            <ul className="divide-y">
              {projects.map((project) => (
                <li key={project.id} className="p-4 flex justify-between items-center">
                  <Link to={`/projects/${project.id}`} className="text-blue-600 hover:underline">
                    {project.name}
                  </Link>
                  <button
                    onClick={() => handleDelete(project.id)}
                    className="text-red-600 hover:text-red-800"
                  >
                    Delete
                  </button>
                </li>
              ))}
            </ul>
          )}
        </div>
      </div>
    </div>
  )
}

