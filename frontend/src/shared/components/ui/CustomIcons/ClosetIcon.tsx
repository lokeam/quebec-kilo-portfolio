import { type SVGProps } from "react";

type ClosetIconProps = SVGProps<SVGSVGElement> & {
  size?: number;
};

const ClosetIcon = ({ size = 24, ...props }: ClosetIconProps) => (
  <svg
    width={size}
    height={size}
    viewBox="0 0 512 512"
    xmlns="http://www.w3.org/2000/svg"
    fill="currentColor"
    {...props}
  >
    <path d="M490.259 0.544H21.741C9.783 0.544 0 10.327 0 22.284v430.471c0 11.958 9.783 21.741 21.741 21.741h31.524v15.219c0 11.958 9.783 21.741 21.741 21.741 11.958 0 21.741-9.783 21.741-21.741v-15.219h317.418v15.219c0 11.958 9.783 21.741 21.741 21.741s20.654-9.783 21.741-21.741v-15.219h31.524c11.958 0 21.741-9.783 22.828-21.741V22.284c0-11.957-9.783-21.74-21.741-21.74zm-169.58 42.394v265.24H42.395V42.938h278.284zm-278.284 308.722h278.285v79.355H42.395v-79.355zm426.123 79.355H364.161V42.938h104.357v388.077z" />
  </svg>
);

export { ClosetIcon };