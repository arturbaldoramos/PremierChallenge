import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import { Suspense, lazy } from 'react'
import { DashboardLayout } from '@/components/layout/DashboardLayout'
import Home from '@/pages/Home'

const Dashboard = lazy(() => import('@/pages/Dashboard'))
const Hospitals = lazy(() => import('@/pages/Hospitals'))
const Doctors = lazy(() => import('@/pages/Doctors'))
const States = lazy(() => import('@/pages/States'))
const Patients = lazy(() => import('@/pages/Patients'))
const Cid = lazy(() => import('@/pages/Cid'))
const Municipalities = lazy(() => import('@/pages/Municipalities'))
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
                 path="/states"
                 element={
                   <DashboardLayout>
                     <Suspense fallback={<LoadingFallback />}>
                       <States />
                     </Suspense>
                   </DashboardLayout>
                 }
               />
               <Route
                 path="/patients"
                 element={
                   <DashboardLayout>
                     <Suspense fallback={<LoadingFallback />}>
                       <Patients />
                     </Suspense>
                   </DashboardLayout>
                 }
               />
               <Route
                 path="/cid"
                 element={
                   <DashboardLayout>
                     <Suspense fallback={<LoadingFallback />}>
                       <Cid />
                     </Suspense>
                   </DashboardLayout>
                 }
               />
               <Route
                 path="/municipalities"
                 element={
                   <DashboardLayout>
                     <Suspense fallback={<LoadingFallback />}>
                       <Municipalities />
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
