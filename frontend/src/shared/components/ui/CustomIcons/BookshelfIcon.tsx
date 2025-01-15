import { type SVGProps } from "react";

type BookshelfIconProps = SVGProps<SVGSVGElement> & {
  size?: number;
};

const BookshelfIcon = ({ size = 24, ...props }: BookshelfIconProps) => (
  <svg
    width={size}
    height={size}
    viewBox="0 0 508 508"
    xmlns="http://www.w3.org/2000/svg"
    fill="currentColor"
    {...props}
  >
    <path d="M493.949 0h-479.8c-7.8 0-14.1 6.3-14.1 14.1v464.1c0 7.8 6.3 14.1 14.1 14.1h24v1.6c0 7.8 6.3 14.1 14.1 14.1s14.1-6.3 14.1-14.1v-1.6h375.3v1.6c0 7.8 6.3 14.1 14.1 14.1s14.1-6.3 14.1-14.1v-1.6h24c7.8 0 14.1-6.3 14.1-14.1V14.1c.1-7.8-6.2-14.1-14-14.1zm-14.1 464.1h-451.6V337.6h451.6v126.5zm-303.2-154.7v-85.5l29.8 85.5h-29.8zm303.2 0h-243.5l-34.5-99.2-25.1 8.7v-7.2h-28.2v97.7h-15.6v-97.7h-28.2v97.7h-15.8v-97.7h-28.2v97.7h-32.5V182.9h451.6v126.5zm0-154.7h-347V57h-28.2v97.7h-15.7V57h-28.2v97.7h-32.5V28.2h451.6v126.5z" />
  </svg>
);

export { BookshelfIcon };
