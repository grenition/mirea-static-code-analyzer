import { Link } from 'react-router-dom'

export default function Home() {
  return (
    <div className="container mx-auto px-4 py-8">
      <div className="max-w-4xl mx-auto">
        <h1 className="text-4xl font-bold mb-6">Static Code Analyzer</h1>
        <p className="text-lg mb-8 text-gray-700">
          Analyze your code in multiple languages: Python, Java, JavaScript, C#, C++, and JSON.
        </p>
        
        <div className="grid md:grid-cols-2 gap-6 mb-8">
          <div className="bg-white p-6 rounded-lg shadow">
            <h2 className="text-2xl font-semibold mb-4">Projects</h2>
            <p className="mb-4 text-gray-600">
              Create and manage projects with multiple files. Upload ZIP archives or create files manually.
            </p>
            <Link to="/projects" className="text-blue-600 hover:underline">
              Go to Projects →
            </Link>
          </div>
          
          <div className="bg-white p-6 rounded-lg shadow">
            <h2 className="text-2xl font-semibold mb-4">Sandbox</h2>
            <p className="mb-4 text-gray-600">
              Analyze a single file without creating a project. Perfect for quick checks.
            </p>
            <Link to="/sandbox" className="text-blue-600 hover:underline">
              Go to Sandbox →
            </Link>
          </div>
        </div>

        <div className="bg-blue-50 p-6 rounded-lg">
          <h2 className="text-xl font-semibold mb-3">Features</h2>
          <ul className="list-disc list-inside space-y-2 text-gray-700">
            <li>Support for 6 programming languages</li>
            <li>Real-time code analysis</li>
            <li>Project management with file tree</li>
            <li>ZIP archive upload support</li>
            <li>Line-by-line issue reporting</li>
          </ul>
        </div>
      </div>
    </div>
  )
}

