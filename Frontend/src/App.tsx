import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import { Suspense, lazy } from 'react'
import { DashboardLayout } from '@/components/layout/DashboardLayout'
import Home from '@/pages/Home'

// Lazy loading do Dashboard
const Dashboard = lazy(() => import('@/pages/Dashboard'))

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route 
          path="/dashboard" 
          element={
            <DashboardLayout>
              <Suspense fallback={
                <div className="flex items-center justify-center h-full">
                  <div className="text-center">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-2"></div>
                    <p className="text-muted-foreground">Carregando dashboard...</p>
                  </div>
                </div>
              }>
                <Dashboard />
              </Suspense>
            </DashboardLayout>
          } 
        />
      </Routes>
    </Router>
  )
}

export default App
