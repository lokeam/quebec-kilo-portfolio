import { Gamepad2, Clapperboard, CircleDollarSign, Package } from 'lucide-react';
import type { LucideIcon } from 'lucide-react';
import type { StatCardIconType } from '../types/singleStatCard.types';

export const SINGLE_STATISTICS_ICONS: Record<StatCardIconType, LucideIcon> = {
  'games': Gamepad2,
  'movies': Clapperboard,
  'onlineServices': CircleDollarSign,
  'package': Package,
  'coin': CircleDollarSign,
};

