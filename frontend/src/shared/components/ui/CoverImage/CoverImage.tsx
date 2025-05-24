import { cn } from "@/shared/components/ui/utils";

type CoverProps = {
  /** The raw `cover.url` you get back from IGDB, e.g. "//images.igdb.com/…/t_thumb/…jpg" */
  src: string;
  /** Any of the IGDB size tokens: "cover_small", "cover_big", "thumb", etc. */
  size?: string;
  alt?: string;
  className?: string;
}

export function CoverImage({ src, size = 'cover_big', alt, className }: CoverProps) {
  // 1) prefix the protocol (IGDB omits http(s):)
  // 2) swap out the size token in the path
  const url = `https:${src}`.replace(/t_[^/]+/, `t_${size}`);

  return (
    <img
      src={url}
      alt={alt ?? 'Game cover'}
      className={cn("object-cover", className)}
    />
  );
}