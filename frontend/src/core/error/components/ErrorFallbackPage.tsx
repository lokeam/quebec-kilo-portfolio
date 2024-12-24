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
  resetErrorBoundary,
  config
}: ErrorFallbackPageProps) => {

  // Log error for monitoring
  console.error('[ErrorBoundary]:', error);

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
            onClick={config?.onAction || resetErrorBoundary}
            label={config?.actionLabel}
          />
        </Box>
      </StyledPaper>
    </Container>
  );
};