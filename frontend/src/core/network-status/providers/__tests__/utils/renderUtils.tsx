import { render } from '@testing-library/react';
import { NetworkStatusProvider } from '../../NetworkStatusProvider';
import { ReactNode } from 'react';

export const renderWithProvider = (component: ReactNode) => {
  return render(
    <NetworkStatusProvider>
      {component}
    </NetworkStatusProvider>
  );
};