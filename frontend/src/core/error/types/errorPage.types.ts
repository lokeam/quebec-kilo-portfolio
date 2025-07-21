import React from 'react';

export type ErrorPageVariant = 'oneCol' | 'twoCol';

export interface ErrorPageBaseProps {
  title: string;
  subtext?: string;
  buttonText?: string;
  onButtonClick?: () => void;
  children?: React.ReactNode;
}

export interface ErrorPageOneColTemplateProps extends ErrorPageBaseProps {
  image?: React.ReactNode;
}

export interface ErrorPageTwoColTemplateProps extends ErrorPageBaseProps {
  image?: React.ReactNode;
}

export interface ErrorPageProps extends ErrorPageBaseProps {
  variant: ErrorPageVariant;
  image?: React.ReactNode;
}