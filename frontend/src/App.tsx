import { BrowserRouter, Route, Routes } from 'react-router-dom';
import HomePage from '@/pages/HomePage/HomePage';
import { ErrorBoundaryProvider } from '@/core/error/providers/ErrorBoundaryProvider';
import { ThemeProvider } from '@/core/theme/providers/ThemeProvider';

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
      <ThemeProvider enableSystemPreference>
        <BrowserRouter>
          <Routes>
            <Route path="/" element={<HomePage />} />
          </Routes>
        </BrowserRouter>
      </ThemeProvider>
    </ErrorBoundaryProvider>
   )
}

export default App;
