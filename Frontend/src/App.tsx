import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import { Suspense, lazy } from 'react'
import { DashboardLayout } from '@/components/layout/DashboardLayout'

const Dashboard = lazy(() => import('@/pages/Dashboard'))
const FileManager = lazy(() => import('@/pages/FileManager'))

function App() {
  const LoadingFallback = () => (
    <div className="flex items-center justify-center h-full">
      <div className="text-center">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-2"></div>
        <p className="text-muted-foreground">Carregando...</p>
      </div>
    </div>
  )

  return (
    <Router>
      <Routes>
        <Route 
          path="/" 
          element={
            <DashboardLayout>
              <Suspense fallback={<LoadingFallback />}>
                <Dashboard />
              </Suspense>
            </DashboardLayout>
          } 
        />
        <Route 
          path="/dashboard" 
          element={
            <DashboardLayout>
              <Suspense fallback={<LoadingFallback />}>
                <Dashboard />
              </Suspense>
            </DashboardLayout>
          } 
        />
        <Route
          path="/file-manager"
          element={
            <DashboardLayout>
              <Suspense fallback={<LoadingFallback />}>
                <FileManager />
              </Suspense>
            </DashboardLayout>
          }
        />
      </Routes>
    </Router>
  )
}

export default App
