import { Gamepad2, Clapperboard, CircleDollarSign, Package } from 'lucide-react';
import { IconCloudDataConnection } from '@tabler/icons-react';
import { type Icon } from '@tabler/icons-react';
import { type LucideIcon, Building } from 'lucide-react';
import type { StatCardIconType } from '../types/singleStatCard.types';

export const SINGLE_STATISTICS_ICONS: Record<StatCardIconType, LucideIcon | Icon> = {
  'games': Gamepad2,
  'movies': Clapperboard,
  'onlineServices': IconCloudDataConnection,
  'package': Package,
  'location': Building,
  'coin': CircleDollarSign,
};

