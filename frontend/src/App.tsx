import { BrowserRouter, Route, Routes } from 'react-router-dom';
import HomePage from '@/pages/HomePage/HomePage';
import { ErrorBoundaryProvider } from '@/core/error/providers/ErrorBoundaryProvider';

import './App.css'

function App() {
  return (
    <ErrorBoundaryProvider
      config={{
        severity: 'error',
        message: 'Application Error',
        actionLabel: 'Reload',
      }}
    >
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<HomePage />} />
        </Routes>
      </BrowserRouter>
    </ErrorBoundaryProvider>
   )
}

export default App
