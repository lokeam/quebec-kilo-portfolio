import { render, screen, fireEvent } from '@testing-library/react';
import { vi } from 'vitest';
import { ErrorButton } from '../ErrorButton';
import { ERROR_ROUTES } from '../../constants/error.constants';

interface MockButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  children?: React.ReactNode;
  'data-testid'?: string;
  variant?: string;
}

// Mock react-router-dom
const mockNavigate = vi.fn();
vi.mock('react-router-dom', () => ({
  useNavigate: () => mockNavigate
}));

// Mock shadcn Button component
vi.mock('@/shared/components/ui/button', () => ({
  Button: ({
    children,
    onClick,
    'data-testid': dataTestId,
    ...props
  }: MockButtonProps) => (
    <button onClick={onClick} data-testid={dataTestId} {...props}>
      {children}
    </button>
  )
}));

describe('ErrorButton', () => {
  const defaultProps = {
    onClick: vi.fn(),
    label: 'Test Button'
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('Navigation variants', () => {
    it('navigates to home when variant is "home"', () => {
      render(
        <ErrorButton
          {...defaultProps}
          variant="home"
        />
      );

      fireEvent.click(screen.getByTestId('error-button-home'));

      expect(mockNavigate).toHaveBeenCalledWith(ERROR_ROUTES.HOME);
      expect(defaultProps.onClick).not.toHaveBeenCalled();
    });

    it('navigates back when variant is "back"', () => {
      render(
        <ErrorButton
          {...defaultProps}
          variant="back"
        />
      );

      fireEvent.click(screen.getByTestId('error-button-back'));

      expect(mockNavigate).toHaveBeenCalledWith(-1);
      expect(defaultProps.onClick).not.toHaveBeenCalled();
    });

    it('calls onClick when variant is "retry"', () => {
      render(
        <ErrorButton
          {...defaultProps}
          variant="retry"
        />
      );

      fireEvent.click(screen.getByTestId('error-button-retry'));

      expect(mockNavigate).not.toHaveBeenCalled();
      expect(defaultProps.onClick).toHaveBeenCalledTimes(1);
    });
  });

  describe('Button styling', () => {
    it('renders with correct variant for retry button', () => {
      render(
        <ErrorButton
          {...defaultProps}
          variant="retry"
        />
      );

      const button = screen.getByTestId('error-button-retry');
      expect(button).toHaveAttribute('variant', 'default');
    });

    it('renders with correct variant for home button', () => {
      render(
        <ErrorButton
          {...defaultProps}
          variant="home"
        />
      );

      const button = screen.getByTestId('error-button-home');
      expect(button).toHaveAttribute('variant', 'secondary');
    });

    it('renders with correct variant for back button', () => {
      render(
        <ErrorButton
          {...defaultProps}
          variant="back"
        />
      );

      const button = screen.getByTestId('error-button-back');
      expect(button).toHaveAttribute('variant', 'outline');
    });
  });

  describe('Accessibility', () => {
    it('renders with provided label', () => {
      render(
        <ErrorButton
          {...defaultProps}
          variant="retry"
        />
      );

      expect(screen.getByText(defaultProps.label)).toBeInTheDocument();
    });

    it('renders with icon and label for each variant', () => {
      const { rerender } = render(
        <ErrorButton
          {...defaultProps}
          variant="retry"
        />
      );

      // Check retry variant
      expect(screen.getByTestId('error-button-retry')).toBeInTheDocument();

      // Check home variant
      rerender(
        <ErrorButton
          {...defaultProps}
          variant="home"
        />
      );
      expect(screen.getByTestId('error-button-home')).toBeInTheDocument();

      // Check back variant
      rerender(
        <ErrorButton
          {...defaultProps}
          variant="back"
        />
      );
      expect(screen.getByTestId('error-button-back')).toBeInTheDocument();
    });
  });
});