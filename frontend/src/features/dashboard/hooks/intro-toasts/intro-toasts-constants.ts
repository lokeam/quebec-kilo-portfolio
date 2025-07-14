export interface IntroToast {
  id: number;
  title: string;
  message: string;
}

export const INTRO_TOASTS: IntroToast[] = [
  {
    id: 1,
    title: "Adding Games to Library",
    message: "Before you may store games you need a place to keep them. Visit the <bold>online-services</bold> or <bold>physical locations</bold> page to add a location.",
  },
  {
    id: 2,
    title: "Storing Physical Games",
    message: "Click <bold>Add Physical Location</bold> to create a place to store your games.",
  },
  {
    id: 3,
    title: "Storing Physical Games",
    message: "Now that you have a physical location, click <bold>Add Sublocation</bold> to create storage spaces for your games.",
  },
  {
    id: 4,
    title: "Storing Digital Games",
    message: "Click <bold>Add Digital Location</bold> to create a place to store your digital games.",
  },
  {
    id: 5,
    title: "Track spending on sublocations",
    message: "QKO automatically tracks digital location <bold>subscription costs</bold>.",
  },
  {
    id: 6,
    title: "Track spending on one time purchases",
    message: "One time purchase tracking helps you remember money spent on games, in-game-purchases DLCs, hardware and other items.",
  },
  {
    id: 7,
    title: "Checking where games are stored",
    message: "The list view button will show details about where you stored your games.",
  },
];