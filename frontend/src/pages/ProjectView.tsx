import { useEffect, useState } from 'react'
import { useParams, useNavigate, Link } from 'react-router-dom'
import axios from 'axios'

interface File {
  id: number
  path: string
  content: string
}

interface Project {
  id: number
  name: string
  files: File[]
}

export default function ProjectView() {
  const { id, file: filePath } = useParams()
  const [project, setProject] = useState<Project | null>(null)
  const [selectedFile, setSelectedFile] = useState<File | null>(null)
  const [analysis, setAnalysis] = useState<any>(null)
  const [analyzerType, setAnalyzerType] = useState('python')
  const navigate = useNavigate()

  useEffect(() => {
    loadProject()
  }, [id])

  useEffect(() => {
    if (filePath && project) {
      const file = project.files.find((f) => f.path === filePath)
      if (file) {
        setSelectedFile(file)
        analyzeFile(file)
      }
    }
  }, [filePath, project])

  const loadProject = async () => {
    try {
      const token = localStorage.getItem('token')
      const response = await axios.get(`/api/projects/${id}`, {
        headers: { Authorization: token },
      })
      setProject(response.data)
    } catch (err: any) {
      if (err.response?.status === 401) {
        navigate('/authorization')
      }
    }
  }

  const analyzeFile = async (file: File) => {
    try {
      const token = localStorage.getItem('token')
      const response = await axios.post(
        `/api/projects/${id}/files/${file.id}/analyze?analyzer=${analyzerType}`,
        {},
        { headers: { Authorization: token } }
      )
      setAnalysis(response.data)
    } catch (err) {
      console.error('Analysis failed:', err)
    }
  }

  const handleFileClick = (file: File) => {
    navigate(`/projects/${id}/${file.path}`)
  }

  if (!project) return <div>Loading...</div>

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="mb-4">
        <Link to="/projects" className="text-blue-600 hover:underline">← Back to Projects</Link>
        <h1 className="text-3xl font-bold mt-2">{project.name}</h1>
      </div>

      <div className="grid grid-cols-3 gap-4">
        <div className="col-span-1 bg-white rounded-lg shadow p-4">
          <h2 className="font-semibold mb-2">Files</h2>
          <ul className="space-y-1">
            {project.files.map((file) => (
              <li key={file.id}>
                <button
                  onClick={() => handleFileClick(file)}
                  className={`text-left w-full px-2 py-1 rounded ${
                    selectedFile?.id === file.id ? 'bg-blue-100' : 'hover:bg-gray-100'
                  }`}
                >
                  {file.path}
                </button>
              </li>
            ))}
          </ul>
        </div>

        <div className="col-span-2">
          {selectedFile ? (
            <div className="bg-white rounded-lg shadow p-4">
              <div className="mb-4 flex items-center gap-2">
                <select
                  value={analyzerType}
                  onChange={(e) => setAnalyzerType(e.target.value)}
                  className="px-3 py-1 border rounded"
                >
                  <option value="python">Python</option>
                  <option value="javascript">JavaScript</option>
                  <option value="java">Java</option>
                  <option value="cpp">C++</option>
                  <option value="csharp">C#</option>
                  <option value="json">JSON</option>
                </select>
                <button
                  onClick={() => analyzeFile(selectedFile)}
                  className="bg-blue-600 text-white px-4 py-1 rounded hover:bg-blue-700"
                >
                  Analyze
                </button>
              </div>

              <div className="mb-4">
                <h3 className="font-semibold mb-2">Code</h3>
                <pre className="bg-gray-100 p-4 rounded overflow-auto max-h-96">
                  <code>{selectedFile.content}</code>
                </pre>
              </div>

              {analysis && (
                <div>
                  <h3 className="font-semibold mb-2">Analysis</h3>
                  <div className="bg-gray-50 p-4 rounded">
                    <pre className="whitespace-pre-wrap">{JSON.stringify(analysis, null, 2)}</pre>
                  </div>
                </div>
              )}
            </div>
          ) : (
            <div className="bg-white rounded-lg shadow p-8 text-center text-gray-500">
              Select a file to view and analyze
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

