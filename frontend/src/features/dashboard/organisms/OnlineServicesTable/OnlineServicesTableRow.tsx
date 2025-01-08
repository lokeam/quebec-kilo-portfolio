import { memo } from 'react';
import { TableCell, TableRow } from "@/shared/components/ui/table";
import { Checkbox } from "@/shared/components/ui/checkbox";
import { Switch } from "@/shared/components/ui/switch";
import { Monitor } from 'lucide-react';
import SVGLogo from "@/shared/components/ui/LogoMap/LogoMap";
import type{ LogoName } from "@/shared/components/ui/LogoMap/LogoMap";
import { PaymentIcon } from 'react-svg-credit-card-payment-icons/dist';
import type { OnlineService } from '@/features/dashboard/pages/OnlineServices/onlineServicesPage.mockdata';
import type { PaymentMethodType } from '@/features/dashboard/pages/OnlineServices/onlineServicesPage.mockdata';

interface OnlineServicesTableRowProps {
  service: OnlineService;
  index: number
};

function OnlineServicesTableRowComponent({ service }: OnlineServicesTableRowProps) {
  const hasValidLogo = Boolean(service.logo);
  const paymentDate = `${service.renewalMonth} ${service.renewalDay}`;
  const isFree = service.billingCycle === 'NA';
  console.log(`${service.label} ${service.billingCycle}`);

  return (
    <TableRow className="h-[72px]">
      <TableCell>
        <Checkbox />
      </TableCell>
      <TableCell>
        <div className="flex items-center gap-3">
          <div className="h-12 w-12 rounded-lg bg-black flex items-center justify-center">
          {hasValidLogo ? (
              <SVGLogo
                domain="games"
                name={service.logo as LogoName<'games'>}
                className="h-8 w-8"
              />
            ) : (
              <Monitor className="h-8 w-8 text-slate-500" />
            )}
          </div>
          <div className="flex flex-col">
            <span className="font-medium">{service.label}</span>
            <span className="text-sm text-muted-foreground">
              {service.tierName || "Standard subscription"}
            </span>
          </div>
        </div>
      </TableCell>
      <TableCell>
        {
          service.isActive ? (
            <Switch checked={true} />
          ) : (
            <Switch checked={false} />
          )
        }
      </TableCell>
      <TableCell>
        {
          isFree ? (
            <span>--</span>
          ) : (
            <span>{service.billingCycle}</span>
          )
        }
      </TableCell>
      <TableCell>{service.monthlyFee}</TableCell>
      <TableCell>
        <PaymentIcon
          type={(service.paymentMethod || 'Generic') as PaymentMethodType}
          format="flatRounded"
        />
      </TableCell>
      <TableCell>
        {
          isFree ? (
            <span>--</span>
          ) : (
            <span>{paymentDate}</span>
          )
        }
      </TableCell>
    </TableRow>
  );
}

export const OnlineServicesTableRow = memo(
  OnlineServicesTableRowComponent,
  (prevProps, nextProps) => {
    return (
      prevProps.index === nextProps.index &&
      prevProps.service.name === nextProps.service.name &&
      prevProps.service.label === nextProps.service.label &&
      prevProps.service.logo === nextProps.service.logo &&
      prevProps.service.tierName === nextProps.service.tierName &&
      prevProps.service.billingCycle === nextProps.service.billingCycle &&
      prevProps.service.monthlyFee === nextProps.service.monthlyFee &&
      prevProps.service.paymentMethod === nextProps.service.paymentMethod
    );
  }
);
