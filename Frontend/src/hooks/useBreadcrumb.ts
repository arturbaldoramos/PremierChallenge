import { useState, useEffect } from 'react'
import { useLocation } from 'react-router-dom'

export interface BreadcrumbItem {
  label: string
  href: string
  icon?: React.ComponentType<{ className?: string }>
}

const breadcrumbMap: Record<string, BreadcrumbItem[]> = {
  '/': [
    { label: 'Dashboard', href: '/dashboard' }
  ],
  '/dashboard': [
    { label: 'Dashboard', href: '/dashboard' }
  ],
  '/file-manager': [
    { label: 'Dashboard', href: '/dashboard' },
    { label: 'Gerenciador de Arquivos', href: '/file-manager' }
  ]
}

export function useBreadcrumb() {
  const location = useLocation()
  const [breadcrumbs, setBreadcrumbs] = useState<BreadcrumbItem[]>([])

  useEffect(() => {
    const currentBreadcrumbs = breadcrumbMap[location.pathname] || [
      { label: 'Dashboard', href: '/dashboard' }
    ]
    setBreadcrumbs(currentBreadcrumbs)
  }, [location.pathname])

  const addBreadcrumb = (item: BreadcrumbItem) => {
    setBreadcrumbs(prev => [...prev, item])
  }

  const removeBreadcrumb = (index: number) => {
    setBreadcrumbs(prev => prev.slice(0, index + 1))
  }

  return {
    breadcrumbs,
    addBreadcrumb,
    removeBreadcrumb
  }
}
