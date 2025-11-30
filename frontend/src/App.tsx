import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import Navbar from './components/Navbar'
import Registration from './pages/Registration'
import Authorization from './pages/Authorization'
import Home from './pages/Home'
import Projects from './pages/Projects'
import ProjectView from './pages/ProjectView'
import Sandbox from './pages/Sandbox'

function App() {
  return (
    <BrowserRouter>
      <div className="min-h-screen bg-gray-50">
        <Navbar />
        <Routes>
          <Route path="/registration" element={<Registration />} />
          <Route path="/authorization" element={<Authorization />} />
          <Route path="/home" element={<Home />} />
          <Route path="/projects" element={<Projects />} />
          <Route path="/projects/:id" element={<ProjectView />} />
          <Route path="/projects/:id/:file" element={<ProjectView />} />
          <Route path="/sandbox" element={<Sandbox />} />
          <Route path="/" element={<Navigate to="/home" replace />} />
        </Routes>
      </div>
    </BrowserRouter>
  )
}

export default App

