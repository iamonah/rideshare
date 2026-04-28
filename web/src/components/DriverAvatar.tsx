import { useState } from "react";

const FALLBACK_DRIVER_AVATAR = `data:image/svg+xml;utf8,${encodeURIComponent(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 64">
  <rect width="64" height="64" rx="32" fill="#1f2937"/>
  <circle cx="32" cy="24" r="12" fill="#f9fafb"/>
  <path d="M12 56c4-11 14-18 20-18s16 7 20 18" fill="#f9fafb"/>
</svg>
`)}`;

interface DriverAvatarProps {
  src?: string;
  alt: string;
  width: number;
  height: number;
  className?: string;
}

export function DriverAvatar({ src, alt, width, height, className }: DriverAvatarProps) {
  const [currentSrc, setCurrentSrc] = useState(src || FALLBACK_DRIVER_AVATAR);

  return (
    <img
      src={currentSrc}
      alt={alt}
      width={width}
      height={height}
      className={className}
      onError={() => setCurrentSrc(FALLBACK_DRIVER_AVATAR)}
    />
  );
}
