export type Message = {
  message: string;
  direction: string
  timestamp: Date;
}

export type User = {
  id: string;
  username: string;
}

export type Product = {
  title: string;
  reviews: {
    username: string;
    rating: number;
    review: string;
    timestamp: Date;
  }[];
}
