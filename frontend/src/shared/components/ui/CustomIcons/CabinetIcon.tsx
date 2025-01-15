import { type SVGProps } from "react";

type CabinetIconProps = SVGProps<SVGSVGElement> & {
  size?: number;
};

const CabinetIcon = ({ size = 24, ...props }: CabinetIconProps) => (
  <svg
    width={size}
    height={size}
    viewBox="3 2 20 20"
    xmlns="http://www.w3.org/2000/svg"
    stroke="currentColor"
    fill="none"
    {...props}
  >
    <rect
      width="16"
      height="16"
      x="4"
      y="3"
      fill="currentColor"
      fillOpacity={0.1}
      rx="1"
    />
    <path
      strokeLinecap="round"
      strokeLinejoin="round"
      strokeWidth="2"
      d="M5 19v2m14-2v2m-3-11v2m-8 0v-2m-3 9h14a1 1 0 0 0 1-1V4a1 1 0 0 0-1-1H5a1 1 0 0 0-1 1v14a1 1 0 0 0 1 1Zm0 0h7V3H5a1 1 0 0 0-1 1v14a1 1 0 0 0 1 1Z"
    />
  </svg>
);

export { CabinetIcon };