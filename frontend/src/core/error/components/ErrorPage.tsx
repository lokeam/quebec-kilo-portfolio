import { ErrorPageOneColTemplate } from '@/core/error/components/ErrorPageOneColTemplate';
import { ErrorPageTwoColTemplate } from '@/core/error/components/ErrorPageTwoColTemplate';

export type ErrorPageVariant = '404' | '500';

export interface ErrorPageProps {
  variant: ErrorPageVariant;
  title: string;
  subtext?: string;
  buttonText?: string;
  role?: 'alert' | 'status';
  ariaLive?: 'polite' | 'assertive';
  onButtonClick?: () => void;
  image?: React.ReactNode;
  children?: React.ReactNode;
}

export const ErrorPage: React.FC<ErrorPageProps> = (props) => {
  const { variant, ...rest } = props;
  if (variant === '404') {
    return <ErrorPageOneColTemplate {...rest} />;
  }
  return <ErrorPageTwoColTemplate {...rest} />;
};