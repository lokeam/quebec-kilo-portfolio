interface LibraryMediaSectionProps {
  title: string
  count: number
  children: React.ReactNode
}

export function LibraryMediaSection({ title, count, children }: LibraryMediaSectionProps) {
  return (
    <>
      <div className="relative w-[calc(100%-45px)] overflow-hidden px-[15px] font-['Julius_Sans_One'] text-base tracking-[2px] text-white">
        {title} <span className="text-[#aaa]">({count})</span>
        <div
          className="absolute ml-[10px] h-[1px] w-full bg-gradient-to-r from-[#404957] to-[#1c2026]"
          style={{ top: '8px' }}
        />
      </div>
      {children}
    </>
  );
}
