import { type SVGProps } from "react";

type DrawerIconProps = SVGProps<SVGSVGElement> & {
  size?: number;
};

const DrawerIcon = ({ size = 24, ...props }: DrawerIconProps) => (
  <svg
    width={size}
    height={size}
    viewBox="0 1 22 22"
    xmlns="http://www.w3.org/2000/svg"
    fill="currentColor"
    {...props}
  >
    <path d="M4 23a1 1 0 0 0 1-1v-1h14v1a1 1 0 0 0 2 0V2a1 1 0 0 0-1-1H4a1 1 0 0 0-1 1v20a1 1 0 0 0 1 1Zm1-11h14v7H5Zm0-9h14v7H5Zm9 12.5a1 1 0 0 1-1 1h-2a1 1 0 0 1 0-2h2a1 1 0 0 1 1 1Zm0-9a1 1 0 0 1-1 1h-2a1 1 0 0 1 0-2h2a1 1 0 0 1 1 1Z" />
  </svg>
);

export { DrawerIcon };