import { Gamepad2, Clapperboard } from 'lucide-react';
import type { LucideIcon } from 'lucide-react';
import type { DomainType } from '@/core/domain/types/domainTypes';

export const DOMAIN_ICONS: Record<DomainType, LucideIcon> = {
  'games': Gamepad2,
  'movies': Clapperboard
};

export const DOMAIN_LABELS = {
  'games': 'Games',
  'movies': 'Movies / Series',
} as const satisfies Record<DomainType, string>;
