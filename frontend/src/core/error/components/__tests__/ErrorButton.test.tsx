import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { ErrorButton } from '../ErrorButton';
import '@testing-library/jest-dom';  // Add this temporarily to verify

describe('ErrorButton', () => {
  // Rendering tests
  it('renders with a default label if none is provided', () => {
    render(<ErrorButton onClick={() => {}} />);
    expect(screen.getByRole('button')).toHaveTextContent('Try Again');
  });

  it('renders with a custom label when provided', () => {
    render(<ErrorButton onClick={() => {}} label='Retry' />);
    expect(screen.getByRole('button')).toHaveTextContent('Retry');
  });

  // Interaction tests
  it('fires the onClick handler when clicked', () => {
    const handleClick = vi.fn();
    render(<ErrorButton onClick={handleClick} />);

    fireEvent.click(screen.getByRole('button'));
    expect(handleClick).toHaveBeenCalledTimes(1);
  });

  // Accessibility tests
  it('has correct ARIA attributes', () => {
    render(<ErrorButton onClick={() => {}} aria-label="Retry action" />);
    const button = screen.getByRole('button');

    expect(button).toHaveAttribute('aria-label', 'Retry action');
  });

  // Style Tests
  it('applies error styles', () => {
    render(<ErrorButton onClick={() => {}} />);
    const button = screen.getByRole('button');

    expect(button).toHaveClass('MuiButton-root');
    expect(button).toHaveClass('MuiButton-contained');
    expect(button).toHaveClass('MuiButton-containedPrimary');
    expect(button).toHaveClass('MuiButton-sizeLarge');
  });

  // Disabled state
  it('respects disabled state', () => {
    render(<ErrorButton onClick={() => {}} disabled />);
    const button = screen.getByRole('button');

    expect(button).toBeDisabled();
    fireEvent.click(button);
  });

});