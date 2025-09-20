import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import { Suspense, lazy } from 'react'
import { DashboardLayout } from '@/components/layout/DashboardLayout'
import Home from '@/pages/Home'

const Dashboard = lazy(() => import('@/pages/Dashboard'))
const Analytics = lazy(() => import('@/pages/Analytics'))
const Users = lazy(() => import('@/pages/Users'))
const Orders = lazy(() => import('@/pages/Orders'))
const Activity = lazy(() => import('@/pages/Activity'))
const Reports = lazy(() => import('@/pages/Reports'))
const Hospitals = lazy(() => import('@/pages/Hospitals'))
const Doctors = lazy(() => import('@/pages/Doctors'))
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
        <Route path="/" element={<Home />} />
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
          path="/analytics" 
          element={
            <DashboardLayout>
              <Suspense fallback={<LoadingFallback />}>
                <Analytics />
              </Suspense>
            </DashboardLayout>
          } 
        />
        <Route 
          path="/users" 
          element={
            <DashboardLayout>
              <Suspense fallback={<LoadingFallback />}>
                <Users />
              </Suspense>
            </DashboardLayout>
          } 
        />
        <Route 
          path="/orders" 
          element={
            <DashboardLayout>
              <Suspense fallback={<LoadingFallback />}>
                <Orders />
              </Suspense>
            </DashboardLayout>
          } 
        />
        <Route 
          path="/activity" 
          element={
            <DashboardLayout>
              <Suspense fallback={<LoadingFallback />}>
                <Activity />
              </Suspense>
            </DashboardLayout>
          } 
        />
        <Route 
          path="/reports" 
          element={
            <DashboardLayout>
              <Suspense fallback={<LoadingFallback />}>
                <Reports />
              </Suspense>
            </DashboardLayout>
          } 
        />
               <Route
                 path="/hospitals"
                 element={
                   <DashboardLayout>
                     <Suspense fallback={<LoadingFallback />}>
                       <Hospitals />
                     </Suspense>
                   </DashboardLayout>
                 }
               />
               <Route
                 path="/doctors"
                 element={
                   <DashboardLayout>
                     <Suspense fallback={<LoadingFallback />}>
                       <Doctors />
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
