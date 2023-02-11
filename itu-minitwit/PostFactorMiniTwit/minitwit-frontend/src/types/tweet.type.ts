import { Author } from "./author.type";

export type Tweet = {
  message_id: number;
  author_id: number;
  text: string;
  pub_date: number;
  flagged: number;
  author: Author;
};
