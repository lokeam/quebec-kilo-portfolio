import { memo, useCallback } from 'react';

// Shadcn Components
import { TableCell, TableRow } from "@/shared/components/ui/table";
import { Checkbox } from "@/shared/components/ui/checkbox";
import { Switch } from "@/shared/components/ui/switch";

// Icons
import SVGLogo from "@/shared/components/ui/LogoMap/LogoMap";
import type { LogoName } from "@/shared/components/ui/LogoMap/LogoMap";
import { Monitor } from 'lucide-react';
import { PaymentIcon } from 'react-svg-credit-card-payment-icons/dist';

// Types
import type { OnlineService } from '@/features/dashboard/lib/types/online-services/services';
import type { PaymentMethod } from '@/shared/constants/payment';

// Hooks
import { toast } from 'sonner';
import { useOnlineServicesToggleActive, useOnlineServicesIsActive } from '@/features/dashboard/lib/stores/onlineServicesStore';
import { isServiceFree } from '@/features/dashboard/lib/utils/online-service-status';

interface OnlineServicesTableRowProps {
  service: OnlineService;
  index: number
};

const createToggleActiveOnlineServiceToast = (label: string, isActive: boolean) => {
  toast(`Recorded ${label} as ${isActive ? 'active' : 'inactive'}`, {
    className: 'bg-green-500 text-white',
    duration: 2500,
  });
};

function OnlineServicesTableRowComponent({ service }: OnlineServicesTableRowProps) {
  const hasValidLogo = Boolean(service.logo);
  const paymentDate = service.billing?.renewalDate ? `${service.billing.renewalDate.month} ${service.billing.renewalDate.day}` : '--';
  const isFree = isServiceFree({ billing: service.billing } as OnlineService);

  // Handlers for online service activation
  const toggleActiveOnlineService = useOnlineServicesToggleActive();
  const isActive = useOnlineServicesIsActive(service.name);

  const handleToggleActiveOnlineService = useCallback((isChecked: boolean) => {
    toggleActiveOnlineService(service.name, isChecked);
    createToggleActiveOnlineServiceToast(service.label, isChecked);
  }, [service.name, service.label, toggleActiveOnlineService]);

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
              {service.tier?.currentTier || "Standard subscription"}
            </span>
          </div>
        </div>
      </TableCell>
      <TableCell>
        <Switch
          checked={isActive}
          onCheckedChange={handleToggleActiveOnlineService}
        />
      </TableCell>
      <TableCell>
        {
          isFree ? (
            <span>--</span>
          ) : (
            <span>{service.billing?.cycle}</span>
          )
        }
      </TableCell>
      <TableCell>{service.billing?.fees.monthly}</TableCell>
      <TableCell>
        <PaymentIcon
          type={(service.billing?.paymentMethod || 'Generic') as PaymentMethod}
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
      prevProps.service.tier?.currentTier === nextProps.service.tier?.currentTier &&
      prevProps.service.billing?.cycle === nextProps.service.billing?.cycle &&
      prevProps.service.billing?.fees.monthly === nextProps.service.billing?.fees.monthly &&
      prevProps.service.billing?.paymentMethod === nextProps.service.billing?.paymentMethod &&
      prevProps.service.billing?.paymentMethod === nextProps.service.billing?.paymentMethod &&
      prevProps.service.status === nextProps.service.status
    );
  }
);
