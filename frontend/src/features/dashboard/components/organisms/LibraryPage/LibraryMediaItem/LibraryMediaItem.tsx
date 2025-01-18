import { cn } from "@/shared/components/ui/utils"
import { Card } from "@/shared/components/ui/card"

interface LibraryMediaItemProps {
  index?: number;
  steamHref?: string
  imageUrl?: string
  className?: string
  favorite?: boolean;
  platformVersion?: string;
}

export function LibraryMediaItem({ steamHref, imageUrl, className }: LibraryMediaItemProps) {

  return (
    <a href={steamHref} className="w-full sm:w-1/2 md:w-1/3 lg:w-1/4 xl:w-1/6 p-2">
      <Card
        className={cn(
          "group rounded-sm relative aspect-[11/15] w-full overflow-hidden border-0 border-t border-t-[#444] shadow-[0_5px_10px_rgba(0,0,0,0.75),_0_0_3px_rgba(123,123,123,0.75)] transition-all duration-250 ease-out hover:shadow-[0_5px_30px_rgba(0,0,0,0.75),_0_0_3px_rgba(123,123,123,0.75)] hover:[transform:perspective(400px)_rotateX(5deg)]",
          className
        )}
        style={{
          backgroundImage: `url(${imageUrl})`,
          backgroundSize: 'cover',
          backgroundPosition: 'center'
        }}
      >
        <div className="card-gradient absolute left-0 top-[-35%] h-full w-full opacity-10 transition-all duration-400 group-hover:top-0 group-hover:opacity-15"
          style={{
            background: 'linear-gradient(30deg, rgba(0, 0, 0, 0), rgba(0, 0, 0, 0) 50%, rgb(255,255,255) 55%)'
          }}
        />
      </Card>
    </a>
  )
}

