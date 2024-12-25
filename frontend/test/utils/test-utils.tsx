import React from 'react';
import { render, RenderOptions } from '@testing-library/react';
import { ErrorBoundaryProvider } from '../../src/core/error/providers/ErrorBoundaryProvider';

export function renderWithProviders(
  ui: React.ReactElement,
  options?: Omit<RenderOptions, 'wrapper'>
) {
  return render(ui, {
    wrapper: ({ children}) => (
      <ErrorBoundaryProvider>
        {children}
      </ErrorBoundaryProvider>
    ),
    ...options,
  });
};
