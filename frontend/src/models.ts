export type Message = {
  message: string;
  sender?: string;
  senderName?: string;
  receiver?: string;
  receiverName?: string;
}

export type Product = {
  title: string;
  reviews: {
    username: string;
    rating: number;
    review: string;
  }[];
}
