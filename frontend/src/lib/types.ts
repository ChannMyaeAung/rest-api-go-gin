export type User = { id: number; email: string; name: string };
export type Event = {
  id: number;
  name: string;
  location: string;
  date: string; // ISO string
  description: string;
  ownerId: number;
};

export type Attendee = {
  id: number;
  userId: number;
  eventId: number;
};

// For displaying attendee lists, we'll use User type
// The backend returns []*User from GetAttendeesByEvent()
export type EventAttendee = User;
