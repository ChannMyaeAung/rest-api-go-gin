export type User = { id: number; email: string; name: string };
export type Event = {
  id: number;
  name: string;
  location: string;
  date: Date; // ISO string
  description: string;
  ownerId: number;
};
export type Attendee = { id: number; name: string; email: string };
