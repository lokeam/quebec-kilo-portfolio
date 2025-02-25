
// Legacy
// export type Game = {
//   id: number;
//   name: string;
//   description: string;
//   coverImage: string;
//   isInLibrary: boolean;
// };

// Updated
export type Game = {
  id: number;
  name: string;
  summary: string;
  rating: number;
  platform_names: string[];
  genre_names: string[];
  theme_names: string[];
  is_in_library: boolean;
  is_in_wishlist: boolean;
}

export const addItemSearchDialogMockData = {
  games: [
    {
      id: 2155,
      name: "Dark Souls",
      summary: "Dark Souls is an action role-playing game developed by FromSoftware and published by Bandai Namco Entertainment...",
      cover_url: "https://images.igdb.com/igdb/image/upload/t_cover_big/co1x78.jpg",
      rating: 87.127762144449,
      platform_names: ["PlayStation 3", "PC (Microsoft Windows)", "Xbox 360"],
      genre_names: ["Role-playing (RPG)", "Adventure"],
      theme_names: ["Action", "Fantasy"],
      is_in_library: false,
      is_in_wishlist: false,
    },
    {
      id: 2368,
      name: "Dark Souls II",
      summary: "Dark Souls II is an action role-playing game developed by FromSoftware and published by Bandai Namco Entertainment...",
      cover_url: "https://images.igdb.com/igdb/image/upload/t_cover_big/co2eoo.jpg",
      rating: 79.10267909602625,
      platform_names: ["PlayStation 3", "PC (Microsoft Windows)", "Xbox 360"],
      genre_names: ["Role-playing (RPG)", "Adventure"],
      theme_names: ["Action", "Fantasy"],
      is_in_library: false,
      is_in_wishlist: false,
    },
    {
      id: 11133,
      name: "Dark Souls III",
      summary: "Dark Souls III is an action role-playing game developed by FromSoftware and published by Bandai Namco Entertainment...",
      cover_url: "https://images.igdb.com/igdb/image/upload/t_cover_big/co2uro.jpg",
      rating: 87.00447070974093,
      platform_names: ["PlayStation 4", "PC (Microsoft Windows)", "Xbox One"],
      genre_names: ["Role-playing (RPG)", "Adventure"],
      theme_names: ["Action", "Fantasy"],
      is_in_library: false,
      is_in_wishlist: false,
    },
    {
      id: 81085,
      name: "Dark Souls: Remastered",
      summary: "Dark Souls Remastered is a remastered version of the original game Dark Souls...",
      cover_url: "https://images.igdb.com/igdb/image/upload/t_cover_big/co2uro.jpg",
      rating: 87.10189857632972,
      platform_names: ["PlayStation 4", "PC (Microsoft Windows)", "Xbox One", "Nintendo Switch"],
      genre_names: ["Role-playing (RPG)", "Adventure"],
      theme_names: ["Action", "Fantasy"],
      is_in_library: false,
      is_in_wishlist: false,
    },
  ],
  total: 4,
}

// Legacy
// export const addItemSearchDialogMockData: Game[] = [
//   {
//     id: 1,
//     name: 'Bloodborne: Game of the Year Edition',
//     description: 'Description 1',
//     coverImage: 'https://images.igdb.com/igdb/image/upload/t_cover_big/co4vui.webp',
//     isInLibrary: true,
//   },
//   {
//     id: 2,
//     name: 'ELDEN RING',
//     description: 'Description 2',
//     coverImage: 'https://images.igdb.com/igdb/image/upload/t_cover_big/co4jni.webp',
//     isInLibrary: true,
//   },
//   {
//     id: 3,
//     name: 'Dark Souls III',
//     description: 'Description 3',
//     coverImage: 'https://images.igdb.com/igdb/image/upload/t_cover_big/co1vcf.webp',
//     isInLibrary: true,
//   },
//   {
//     id: 4,
//     name: 'Helldivers 2',
//     description: 'Description 4',
//     coverImage: 'https://images.igdb.com/igdb/image/upload/t_cover_big/co741o.webp',
//     isInLibrary: true,
//   },
//   {
//     id: 5,
//     name: 'Demon\'s Souls',
//     description: 'Description 5',
//     coverImage: 'https://images.igdb.com/igdb/image/upload/t_cover_big/co2kj9.webp',
//     isInLibrary: true,
//   },
//   {
//     id: 6,
//     name: 'Stardew Valley',
//     description: 'Description 6',
//     coverImage: 'https://images.igdb.com/igdb/image/upload/t_cover_big/xrpmydnu9rpxvxfjkiu7.webp',
//     isInLibrary: true,
//   },
//   {
//     id: 7,
//     name: 'Gran Turismo 7',
//     description: 'Description 6',
//     coverImage: 'https://images.igdb.com/igdb/image/upload/t_cover_big/co2g84.webp',
//     isInLibrary: false,
//   },
//   {
//     id: 8,
//     name: 'Metal Slug Tactics',
//     description: 'Description 7',
//     coverImage: 'https://images.igdb.com/igdb/image/upload/t_cover_big/co8c4a.webp',
//     isInLibrary: false,
//   },
//   {
//     id: 9,
//     name: 'Vagrant Story',
//     description: 'Description 8',
//     coverImage: 'https://images.igdb.com/igdb/image/upload/t_cover_big/co2rso.webp',
//     isInLibrary: false,
//   },
//   {
//     id: 10,
//     name: 'Bust-a-Move',
//     description: 'Description 9',
//     coverImage: 'https://images.igdb.com/igdb/image/upload/t_cover_big/co6bkg.webp',
//     isInLibrary: true,
//   },
//   {
//     id: 10,
//     name: 'Rad Racer',
//     description: 'Description 10',
//     coverImage: 'https://images.igdb.com/igdb/image/upload/t_cover_big/co6049.webp',
//     isInLibrary: false,
//   },
//   {
//     id: 11,
//     name: 'The King of Fighters XV',
//     description: 'Description 11',
//     coverImage: 'https://images.igdb.com/igdb/image/upload/t_cover_big/co6gt4.webp',
//     isInLibrary: false,
//   },
//   {
//     id: 12,
//     name: 'Gradius V',
//     description: 'Description 12',
//     coverImage: 'https://images.igdb.com/igdb/image/upload/t_cover_big/co50f9.webp',
//     isInLibrary: false,
//   },
//   {
//     id: 13,
//     name: 'Gunstar Heroes',
//     description: 'Description 13',
//     coverImage: 'https://images.igdb.com/igdb/image/upload/t_cover_big/co6pmh.webp',
//     isInLibrary: true,
//   },
//   {
//     id: 14,
//     name: 'Age of Empires IV',
//     description: 'Description 14',
//     coverImage: 'https://images.igdb.com/igdb/image/upload/t_cover_big/co39tg.webp',
//     isInLibrary: false,
//   },
//   {
//     id: 15,
//     name: 'Path of Exile 2',
//     description: 'Description 15',
//     coverImage: 'https://images.igdb.com/igdb/image/upload/t_cover_big/co8ae0.webp',
//     isInLibrary: false,
//   },
//   {
//     id: 16,
//     name: 'Chef Chen',
//     description: 'Description 16',
//     coverImage: 'https://images.igdb.com/igdb/image/upload/t_cover_big/co7pix.webp',
//     isInLibrary: false,
//   },
//   {
//     id: 17,
//     name: 'Rad Racer',
//     description: 'Description 17',
//     coverImage: 'https://images.igdb.com/igdb/image/upload/t_cover_big/co6049.webp',
//     isInLibrary: false,
//   },
//   {
//     id: 17,
//     name: 'Beat Saber',
//     description: 'Description 18',
//     coverImage: 'https://images.igdb.com/igdb/image/upload/t_cover_big/co20ux.webp',
//     isInLibrary: false,
//   },
//   {
//     id: 19,
//     name: 'Microsoft Flight Simulator 2024',
//     description: 'Description 19',
//     coverImage: 'https://images.igdb.com/igdb/image/upload/t_cover_big/co91qm.webp',
//     isInLibrary: false,
//   },
//   {
//     id: 20,
//     name: 'Forza Horizon 5',
//     description: 'Description 20',
//     coverImage: 'https://images.igdb.com/igdb/image/upload/t_cover_big/co3ofx.webp',
//     isInLibrary: true,
//   },
//   {
//     id: 21,
//     name: 'The Elder Scrolls V: Skyrim',
//     description: 'Description 21',
//     coverImage: 'https://images.igdb.com/igdb/image/upload/t_cover_big/co1tnw.webp',
//     isInLibrary: false,
//   },
//   {
//     id: 22,
//     name: 'Monster Hunter Rise + Sunbreak',
//     description: 'Description 22',
//     coverImage: 'https://images.igdb.com/igdb/image/upload/t_cover_big/co7i5f.webp',
//     isInLibrary: false,
//   },
// ];
