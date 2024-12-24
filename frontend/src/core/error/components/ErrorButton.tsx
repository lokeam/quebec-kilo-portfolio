import { Button } from '@mui/material';
import { styled } from '@mui/material/styles';

const StyledButton = styled(Button)(({ theme }) => ({
  margin: theme.spacing(1),
  minWidth: 100,
}));

interface ErrorButtonProps {
  onClick: () => void;
  label?: string;
}

export const ErrorButton = ({ onClick, label = 'Try Again' } : ErrorButtonProps) => {
  console.log('ErrorButton');
  return (
    <StyledButton
      variant="contained"
      color="primary"
      onClick={onClick}
      size="large"
    >
      {label}
    </StyledButton>
  );
};


