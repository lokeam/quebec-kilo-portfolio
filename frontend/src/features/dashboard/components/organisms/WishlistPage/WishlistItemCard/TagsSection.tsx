import { memo } from 'react';
import { Badge } from "@/shared/components/ui/badge";

interface TagsSectionProps {
  tags: string[];
}

export const TagsSection = memo(({ tags }: TagsSectionProps) => {
  return (
    <div className="flex flex-nowrap gap-2 overflow-hidden h-8">
      {tags.map((tag, index) => (
        <Badge
          key={index}
          variant="secondary"
          className="bg-[#42464e] hover:bg-[#42464e] rounded-none text-nowrap whitespace-nowrap no-wrap"
        >
          {tag}
        </Badge>
      ))}
    </div>
  );
}, (prevProps, nextProps) => {
  return prevProps.tags.length === nextProps.tags.length &&
    prevProps.tags.every((tag, i) => tag === nextProps.tags[i]);
});
