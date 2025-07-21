import { Card, CardContent } from '@/shared/components/ui/card'
import { Star, MoreHorizontal } from '@/shared/components/ui/icons'
import { IconShoppingCartDollar } from '@/shared/components/ui/icons'
import { Button } from '@/shared/components/ui/button'

type WishListDealsCardProps = {
  starredItem: string;
  starredItemCurrentPrice: string;
  itemsOnSale: number;
  cheapestSaleItemPercentage: string;
};

export function WishListDealsCard({
  starredItem,
  starredItemCurrentPrice,
  itemsOnSale,
  cheapestSaleItemPercentage,
}: WishListDealsCardProps) {
  return (
    <div className="flex flex-col gap-6 min-h-[300px] h-full">
      {/* Total Wishlist Items Card */}
      <Card className="bg-gradient-to-b from-slate-900 to-slate-950 border-slate-800 flex-1">
        <CardContent className="p-6">
          <div className="flex justify-between items-center mb-4">
            <div className="h-12 w-12 rounded-full bg-violet-500/20 flex items-center justify-center">
              <Star className="h-6 w-6 text-violet-500" />
            </div>
            <Button variant="ghost" size="icon" className="text-slate-400">
              <MoreHorizontal className="h-5 w-5" />
            </Button>
          </div>
          <div className="space-y-2">
            <p className="text-sm text-slate-400">Starred Wishlist Item</p>
            <div className="flex items-baseline gap-2">
              <h4 className="text-xl font-semibold text-white">{starredItem}</h4>
            </div>
            <p className="text-sm text-slate-500">
              Best Current Price:
              <span className="text-sm font-medium text-white">{starredItemCurrentPrice}</span>
            </p>
          </div>
        </CardContent>
      </Card>

      {/* Deals Card */}
      <Card className="bg-gradient-to-b from-slate-900 to-slate-950 border-slate-800 flex-1">
        <CardContent className="p-6">
          <div className="flex justify-between items-center mb-4">
            <div className="h-12 w-12 rounded-full bg-violet-500/20 flex items-center justify-center">
              <IconShoppingCartDollar className="h-6 w-6 text-violet-500" />
            </div>
            <Button variant="ghost" size="icon" className="text-slate-400">
              <MoreHorizontal className="h-5 w-5" />
            </Button>
          </div>
          <div className="space-y-2">
            <p className="text-sm text-slate-400">Wishlist Items on Sale</p>
            <div className="flex items-baseline gap-2">
              <h4 className="text-3xl font-semibold text-white">{itemsOnSale}</h4>
              {/* <span className="text-sm font-medium text-orange-500">As of 01/04/2025</span> */}
            </div>
            <p className="text-sm text-green-500">Cheapest Item {cheapestSaleItemPercentage}</p>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
