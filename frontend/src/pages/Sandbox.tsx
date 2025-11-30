import { useState, useRef, useEffect } from 'react'
import axios from 'axios'

export default function Sandbox() {
  const [code, setCode] = useState('')
  const [analyzerType, setAnalyzerType] = useState('python')
  const [analysis, setAnalysis] = useState<any>(null)
  const [isAnalyzing, setIsAnalyzing] = useState(false)
  const debounceTimer = useRef<number | null>(null)

  useEffect(() => {
    if (code && code.trim()) {
      if (debounceTimer.current) {
        window.clearTimeout(debounceTimer.current)
      }
      debounceTimer.current = window.setTimeout(() => {
        analyzeCode()
      }, 1000)
    }

    return () => {
      if (debounceTimer.current) {
        window.clearTimeout(debounceTimer.current)
      }
    }
  }, [code, analyzerType])

  const analyzeCode = async () => {
    if (!code.trim()) return

    setIsAnalyzing(true)
    try {
      const response = await axios.post(`/api/analyzer/${analyzerType}`, {
        files: [
          {
            path: `sandbox.${getExtension(analyzerType)}`,
            content: code,
          },
        ],
      })
      setAnalysis(response.data)
    } catch (err) {
      console.error('Analysis failed:', err)
      setAnalysis({ error: 'Analysis failed' })
    } finally {
      setIsAnalyzing(false)
    }
  }

  const getExtension = (type: string) => {
    const extensions: { [key: string]: string } = {
      python: 'py',
      javascript: 'js',
      java: 'java',
      cpp: 'cpp',
      csharp: 'cs',
      json: 'json',
    }
    return extensions[type] || 'txt'
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="max-w-6xl mx-auto">
        <h1 className="text-3xl font-bold mb-6">Sandbox</h1>

        <div className="mb-4 flex items-center gap-2">
          <label className="font-semibold">Analyzer:</label>
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
          {isAnalyzing && <span className="text-gray-500">Analyzing...</span>}
        </div>

        <div className="grid md:grid-cols-2 gap-4">
          <div className="bg-white rounded-lg shadow p-4">
            <h2 className="font-semibold mb-2">Code</h2>
            <textarea
              value={code}
              onChange={(e) => setCode(e.target.value)}
              className="w-full h-96 p-4 border rounded font-mono text-sm"
              placeholder="Enter your code here..."
            />
          </div>

          <div className="bg-white rounded-lg shadow p-4">
            <h2 className="font-semibold mb-2">Analysis Results</h2>
            <div className="h-96 overflow-auto bg-gray-50 p-4 rounded">
              {analysis ? (
                <pre className="whitespace-pre-wrap text-sm">
                  {JSON.stringify(analysis, null, 2)}
                </pre>
              ) : (
                <p className="text-gray-500">Enter code to see analysis results</p>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

