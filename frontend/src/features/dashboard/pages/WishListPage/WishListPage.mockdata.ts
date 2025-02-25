import type { WishlistGroups } from '@/features/dashboard/lib/types/wishlist/base';

export const wishlistPageMockData: WishlistGroups = {
  pc: [
    {
      id: '1',
      name: 'Fallout 4',
      cover_url: 'https://images.igdb.com/igdb/image/upload/t_cover_big/co1yc6.webp',
      theme_names: ['Open-World', 'Post-Apocalyptic', 'Single-Player', 'Exploration'],
      releaseDate: 'Nov 10, 2015',
      rating: {
        positive: 80,
        negative: 20,
        totalReviews: 100,
      },
      price: {
        original: 19.99,
        discounted: 4.99,
        discountPercentage: 75,
        vendor: 'Steam',
      },
      platform: 'PC',
      platformSupport: {
        hasMacOSVersion: false,
      },
      is_in_library: false,
      is_in_wishlist: true,
    },
  ],
  console: [
    {
      id: '4',
      name: 'Dark Souls III',
      cover_url: 'https://images.igdb.com/igdb/image/upload/t_cover_big/co1vcf.jpg',
      theme_names: ['Role-playing (RPG)', 'Adventure', 'Action', 'Fantasy', 'Single Player', 'Multiplayer'],
      releaseDate: 'Nov 1, 2011',
      rating: {
        positive: 300,
        negative: 10,
        totalReviews: 310,
      },
      price: {
        original: 59.99,
        discounted: 49.99,
        discountPercentage: 15,
        vendor: 'iOS App Store',
      },
      platform: 'Console',
      is_in_library: false,
      is_in_wishlist: true,
    },
    {
      id: '3',
      name: 'Gradius V',
      cover_url: 'https://images.igdb.com/igdb/image/upload/t_cover_big/co50f9.webp',
      theme_names: ['Shooter', 'Action', 'Science fiction', 'Single Player', 'Co-operative', 'Side view'],
      releaseDate: 'Jul 7, 2004',
      rating: {
        positive: 9300,
        negative: 10,
        totalReviews: 9310,
      },
      price: {
        original: 59.99,
        discounted: 55.99,
        discountPercentage: 6.7,
        vendor: 'Best Buy',
      },
      platform: 'Console',
      is_in_library: false,
      is_in_wishlist: true,
    },
  ],
  mobile: [
    {
      id: '2',
      name: 'Balatro',
      cover_url: 'https://images.igdb.com/igdb/image/upload/t_cover_big/co752b.webp',
      theme_names: ['Strategy', 'Indie', 'Card & Board Game', 'Single Player', 'Isometric'],
      releaseDate: 'Nov 1, 2011',
      rating: {
        positive: 300,
        negative: 10,
        totalReviews: 310,
      },
      price: {
        original: 14.99,
        discounted: 10.99,
        discountPercentage: 15,
        vendor: 'iOS App Store',
      },
      platform: 'Mobile',
      platformSupport: {
        hasAndroidVersion: true,
        hasIOSVersion: true,
      },
      is_in_library: false,
      is_in_wishlist: true,
    },
    {
      id: '3',
      name: 'Levelhead',
      cover_url: 'https://images.igdb.com/igdb/image/upload/t_cover_big/co1ks9.webp',
      theme_names: ['Platform', 'Simulator', 'Adventure', 'Single Player', 'Side View'],
      releaseDate: 'Nov 1, 2011',
      rating: {
        positive: 700,
        negative: 440,
        totalReviews: 1140,
      },
      price: {
        original: 6.99,
        discounted: 5.99,
        discountPercentage: 15,
        vendor: 'Google Play Store',
      },
      platform: 'Mobile',
      platformSupport: {
        hasAndroidVersion: true,
        hasIOSVersion: true,
      },
      is_in_library: false,
      is_in_wishlist: true,
    },
  ],
}


/*
Legacy:

export const wishlistPageMockData: WishlistGroups = {
  pc: [
    {
      id: '1',
      title: 'Fallout 4',
      thumbnailUrl: 'https://images.igdb.com/igdb/image/upload/t_cover_big/co1yc6.webp',
      tags: ['Open-World', 'Post-Apocalyptic', 'Single-Player', 'Exploration'],
      releaseDate: 'Nov 10, 2015',
      rating: {
        positive: 80,
        negative: 20,
        totalReviews: 100,
      },
      price: {
        original: 19.99,
        discounted: 4.99,
        discountPercentage: 75,
        vendor: 'Steam',
      },
      platform: 'PC',
      platformSupport: {
        hasMacOSVersion: false,
      },
    },
  ],
  console: [
    {
      id: '4',
      title: 'Dark Souls III',
      thumbnailUrl: 'https://images.igdb.com/igdb/image/upload/t_cover_big/co1vcf.webp',
      tags: ['Role-playing (RPG)', 'Adventure', 'Action', 'Fantasy', 'Single Player', 'Multiplayer'],
      releaseDate: 'Nov 1, 2011',
      rating: {
        positive: 300,
        negative: 10,
        totalReviews: 310,
      },
      price: {
        original: 59.99,
        discounted: 49.99,
        discountPercentage: 15,
        vendor: 'iOS App Store',
      },
      platform: 'Console',
    },
    {
      id: '3',
      title: 'Gradius V',
      thumbnailUrl: 'https://images.igdb.com/igdb/image/upload/t_cover_big/co50f9.webp',
      tags: ['Shooter', 'Action', 'Science fiction', 'Single Player', 'Co-operative', 'Side view'],
      releaseDate: 'Jul 7, 2004',
      rating: {
        positive: 9300,
        negative: 10,
        totalReviews: 9310,
      },
      price: {
        original: 59.99,
        discounted: 55.99,
        discountPercentage: 6.7,
        vendor: 'Best Buy',
      },
      platform: 'Console',
    },
  ],
  mobile: [
    {
      id: '2',
      title: 'Balatro',
      thumbnailUrl: 'https://images.igdb.com/igdb/image/upload/t_cover_big/co752b.webp',
      tags: ['Strategy', 'Indie', 'Card & Board Game', 'Single Player', 'Isometric'],
      releaseDate: 'Nov 1, 2011',
      rating: {
        positive: 300,
        negative: 10,
        totalReviews: 310,
      },
      price: {
        original: 14.99,
        discounted: 10.99,
        discountPercentage: 15,
        vendor: 'iOS App Store',
      },
      platform: 'Mobile',
      platformSupport: {
        hasAndroidVersion: true,
        hasIOSVersion: true,
      },
    },
    {
      id: '3',
      title: 'Levelhead',
      thumbnailUrl: 'https://images.igdb.com/igdb/image/upload/t_cover_big/co1ks9.webp',
      tags: ['Platform', 'Simulator', 'Adventure', 'Single Player', 'Side View'],
      releaseDate: 'Nov 1, 2011',
      rating: {
        positive: 700,
        negative: 440,
        totalReviews: 1140,
      },
      price: {
        original: 6.99,
        discounted: 5.99,
        discountPercentage: 15,
        vendor: 'Google Play Store',
      },
      platform: 'Mobile',
      platformSupport: {
        hasAndroidVersion: true,
        hasIOSVersion: true,
      },
    },
  ],
}
*/