import { Button } from '@mui/material';
import { styled } from '@mui/material/styles';

const StyledButton = styled(Button)(({ theme }) => ({
  margin: theme.spacing(1),
  minWidth: 100,
}));

interface ErrorButtonProps {
  onClick: () => void;
  label?: string;
  disabled?: boolean;
  'aria-label'?: string;
  className?: string;
  onKeyDown?: (event: React.KeyboardEvent) => void;
}

export const ErrorButton = ({
  onClick,
  label = 'Try Again',
  disabled = false,
  'aria-label': ariaLabel,
  className = '',
  onKeyDown,
} : ErrorButtonProps) => {
  console.log('ErrorButton');
  return (
    <StyledButton
      variant="contained"
      color="primary"
      onClick={onClick}
      onKeyDown={onKeyDown}
      disabled={disabled}
      className={className}
      aria-label={ariaLabel}
      size="large"
    >
      {label}
    </StyledButton>
  );
};


