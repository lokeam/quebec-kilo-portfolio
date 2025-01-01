import { useNavigate } from 'react-router-dom';
import { Button } from '@/shared/components/ui/button';
import { RefreshCcw, House, ArrowLeft } from 'lucide-react';
import type { ErrorButtonProps } from '../types/error.types';
import { ERROR_ROUTES } from '../constants/error.constants';

export const ErrorButton = ({ onClick, variant, label }: ErrorButtonProps) => {
  const navigate = useNavigate();

  const handleClick = () => {
    switch (variant) {
      case 'home':
        navigate(ERROR_ROUTES.HOME);
        break;
      case 'back':
        navigate(-1);
        break;
      default:
        onClick();
    }
  };

  const getButtonIcon = () => {
    switch (variant) {
      case 'home':
        return <House className="mr-2 h-4 w-4" />;
      case 'back':
        return <ArrowLeft className="mr-2 h-4 w-4" />;
      case 'retry':
        return <RefreshCcw className="mr-2 h-4 w-4" />;
      default:
        return null;
    }
  };

  const getButtonVariant = () => {
    switch (variant) {
      case 'retry':
        return 'default' as const;
      case 'home':
        return 'secondary' as const;
      case 'back':
        return 'outline' as const;
      default:
        return 'default' as const;
    }
  };

  return (
    <Button
      onClick={handleClick}
      variant={getButtonVariant()}
      className="min-w-[100px] m-4"
      data-testid={`error-button-${variant}`}
    >
      {getButtonIcon()}
      <span className="ml-2">{label}</span>
    </Button>
  );
};
