import { render } from '@testing-library/react';
import { NetworkStatusProvider } from '../../NetworkStatusProvider';
import { type ReactElement } from 'react';

export const renderWithProvider = (component: ReactElement) => {
  return render(
    <NetworkStatusProvider>
      {component}
    </NetworkStatusProvider>
  );
};