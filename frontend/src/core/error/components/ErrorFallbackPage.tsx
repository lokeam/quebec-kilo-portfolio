import { Box, Typography, Container, Paper } from '@mui/material';
import { styled } from '@mui/material/styles';
import { FallbackProps } from 'react-error-boundary';
import { ErrorButton } from '@/core/error/components/ErrorButton';
import { ErrorConfig } from '@/core/error/types/error.types';

const StyledPaper = styled(Paper)(({ theme }) => ({
  padding: theme.spacing(4),
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  marginTop: theme.spacing(8),
}));

interface ErrorFallbackPageProps extends FallbackProps {
  config?: ErrorConfig;
}

export const ErrorFallbackPage = ({
  error,
  resetErrorBoundary, // Provided by react-error-boundary
  config
}: ErrorFallbackPageProps) => {

  const handleClick = () => {
    if (config?.onAction) {
      // If custom action exists, only call that
      config.onAction();
    } else {
      // Otherwise use the default reset behavior
      resetErrorBoundary();
    }
  };

  const handleKeyDown = (event: React.KeyboardEvent) => {
    if (event.key === 'Enter' || event.key === ' ') {
       // Prevent page scroll on space
      event.preventDefault();
      if (config?.onAction) {
        config.onAction();
      } else {
        resetErrorBoundary();
      }
    }
  };
  // Log error for monitoring
  console.error('[ErrorBoundary] Error:', error);

  return (
    <Container maxWidth="sm">
      <StyledPaper elevation={3}>
        <Box textAlign="center">
          <Typography variant="h2" color="error" gutterBottom>
            {config?.severity === 'fatal' ? 'Critical Error' : 'Oops!'}
          </Typography>

          <Typography variant="h5" gutterBottom>
            {config?.message || 'Something went wrong'}
          </Typography>

          <Typography variant="body1" color="textSecondary" paragraph>
            {error?.message || 'We apologize for the inconvenience. Please try again.'}
          </Typography>

          <ErrorButton
            onClick={handleClick}
            onKeyDown={handleKeyDown}
            label={config?.actionLabel || 'Try Again'}
          />
        </Box>
      </StyledPaper>
    </Container>
  );
};