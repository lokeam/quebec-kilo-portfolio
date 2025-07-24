// Centralized icon imports to optimize bundle size
// This file consolidates all icon imports to prevent tree-shaking issues
// Only exports icons that are actually used in the codebase

// Tabler Icons - Only import what we actually use
export {
  // Most frequently used (6+ files)
  IconCloudDataConnection,  // Used in 8+ files
  IconEdit,                 // Used in 6+ files
  IconTrash,               // Used in 6+ files

  // Frequently used (2-5 files)
  IconCar,                 // Used in 3+ files
  IconSettings,            // Used in 2+ files
  IconHeart,               // Used in 2+ files
  IconShoppingCartDollar,  // Used in 2+ files
  IconDeviceGamepad2,      // Used in 2+ files
  IconPackage,             // Used in 2+ files
  IconDisc,                // Used in 2+ files
  IconCpu2,                // Used in 2+ files
  IconSparkles,            // Used in 2+ files
  IconCpu,                 // Used in 2+ files
  IconCloudDown,           // Used in 2+ files
  IconAlertTriangle,       // Used in 2+ files

  // Single use icons
  IconCalendarDollar,      // Used in 1 file
  IconMailOff,             // Used in 1 file
  IconCloudOff,            // Used in 1 file
  IconStar,                // Used in 1 file
  IconStarFilled,          // Used in 1 file
  IconLayoutDashboard,     // Used in 1 file
  IconBook,                // Used in 1 file
  IconBrandAndroid,        // Used in 1 file
  IconBrandApple,          // Used in 1 file
  IconBrandWindowsFilled,  // Used in 1 file
  IconDeviceGamepad,       // Used in 1 file
  IconDeviceMobile,        // Used in 1 file

  // Box multiple icons (used in getLibraryItemCountIcon)
  IconBoxMultiple2,
  IconBoxMultiple3,
  IconBoxMultiple4,
  IconBoxMultiple5,
  IconBoxMultiple6,
  IconBoxMultiple7,
  IconBoxMultiple8,
  IconBoxMultiple9,
} from '@tabler/icons-react';

// Lucide Icons - Only import what we actually use
export {
  // Most frequently used (3+ files)
  Package,                 // Used in 5+ files
  House,                   // Used in 4+ files
  Building2,               // Used in 4+ files
  Warehouse,               // Used in 3+ files
  Building,                // Used in 3+ files
  LayoutGrid,              // Used in 3+ files
  LayoutList,              // Used in 3+ files
  BoxIcon,                 // Used in 3+ files
  CalendarIcon,            // Used in 3+ files

  // Frequently used (2 files)
  Gamepad2,                // Used in 2+ files
  Cpu,                     // Used in 2+ files
  Disc,                    // Used in 2+ files
  Sparkles,                // Used in 2+ files
  Cloud,                   // Used in 2+ files
  ChevronsUpDown,          // Used in 2+ files
  Check,                   // Used in 2+ files
  ArrowUpRight,            // Used in 2+ files
  TrendingDown,            // Used in 2+ files
  TrendingUp,              // Used in 2+ files

  // Single use icons
  Bell,                    // Used in 1 file
  Tag,                     // Used in 1 file
  BarChart,                // Used in 1 file
  SquareCheckBig,          // Used in 1 file
  LogOut,                  // Used in 1 file
  Home,                    // Used in 1 file
  LibraryBig,              // Used in 1 file
  Clapperboard,            // Used in 1 file
  CircleDollarSign,        // Used in 1 file
  MoreHorizontal,          // Used in 1 file
  Star,                    // Used in 1 file
  MessageSquare,           // Used in 1 file
  Send,                    // Used in 1 file
  ChevronRight,            // Used in 1 file
  ChevronDown,             // Used in 1 file
  Image,                   // Used in 1 file
  AlertTriangle,           // Used in 1 file
  Mail,                    // Used in 1 file
  ExternalLink,            // Used in 1 file
  RefreshCcw,              // Used in 1 file
  ArrowLeft,               // Used in 1 file
  SearchIcon,              // Used in 1 file
  CheckCircle2,            // Used in 1 file
  CheckCircle,             // Used in 1 file
  Ban,                     // Used in 1 file
  PackagePlus,             // Used in 1 file
  Settings,                // Used in 1 file
  Trash2,                  // Used in 1 file
  Monitor,                 // Used in 1 file
  Moon,                    // Used in 1 file
  Sun,                     // Used in 1 file
  X,                       // Used in 1 file
  Info,                    // Used in 1 file
  Sheet,                   // Used in 1 file
  SquareDashed,            // Used in 1 file
} from 'lucide-react';