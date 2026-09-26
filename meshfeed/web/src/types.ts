export interface FeedView {
  id: string;
  latest_seq: number;
}

export interface ItemView {
  seq: number;
  timestamp: number;
  content: string;
}
