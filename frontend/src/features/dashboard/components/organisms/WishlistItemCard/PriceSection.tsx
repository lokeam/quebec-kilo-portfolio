import { memo } from 'react';
import { Button } from "@/shared/components/ui/button";

interface PriceSectionProps {
  price: {
    original: number;
    discounted?: number | null;
    discountPercentage?: number | null;
    vendor: string;
  };
  showMoreDeals: boolean;
  stackPriceContent: boolean;
}

export const PriceSection = memo(({ price, showMoreDeals, stackPriceContent }: PriceSectionProps) => {
  return (
    <div className={`flex items-center gap-4 ${stackPriceContent ? 'flex-col' : ''}`}>
      <div className="flex items-center gap-2">
        {price.discountPercentage && (
          <div className="bg-[#94d933] text-black font-bold hover:bg-[#567b27] rounded-sm py-2 px-3 mr-3">
            -{price.discountPercentage}%
          </div>
        )}
        <div className="flex flex-col">
          <div className="text-gray-400 line-through text-lg">
            ${price.original.toFixed(2)}
          </div>
        </div>
      </div>
      <div className="flex flex-col gap-2">
        <Button className="bg-[#4c6b22] hover:bg-[#567b27] text-white gap-1 text-wrap">
          <span className="text-md font-bold">
            ${(price.discounted ?? price.original).toFixed(2)}
          </span>
          <span className="text-sm">from {price.vendor}</span>
        </Button>
        {showMoreDeals && (
          <Button className="bg-[#492fef] hover:bg-[#632de1] text-white">
            See more deals
          </Button>
        )}
      </div>
    </div>
  );
}, (prevProps, nextProps) => {
  return (
    prevProps.price.original === nextProps.price.original &&
    prevProps.price.discounted === nextProps.price.discounted &&
    prevProps.price.discountPercentage === nextProps.price.discountPercentage &&
    prevProps.price.vendor === nextProps.price.vendor &&
    prevProps.showMoreDeals === nextProps.showMoreDeals &&
    prevProps.stackPriceContent === nextProps.stackPriceContent
  );
});
