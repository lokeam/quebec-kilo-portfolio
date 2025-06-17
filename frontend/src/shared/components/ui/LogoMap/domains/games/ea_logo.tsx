export default function EALogo({ className }: { className?: string }) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      xmlSpace="preserve"
      width="100px"
      height="100px"
      viewBox="-200 -200 1400 1400" /* Added padding around the original 1000x1000 viewBox */
      className={className}
    >
      <path
        fill="#ff4747"
        fillRule="evenodd"
        d="M500 1000.001c-275.703 0-500-224.3-500-500 0-275.703 224.297-500 500-500 275.698 0 500 224.299 500 500 0 275.7-224.3 500-500 500zm84.628-693.396H302.054l-42.874 68.901h282.25zm57.747.658L469.632 582.325H278.018l44.207-68.96H437.07l43.87-68.925H215.439l-43.861 68.926h62.896l-87.266 137.68h364.19L645.912 438.91l49.046 74.456h-44.224l-41.882 68.959h130.961l45.475 68.721h83.54z"
        clipRule="evenodd"
      />
    </svg>
  );
}
