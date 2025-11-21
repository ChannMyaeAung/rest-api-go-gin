"use client";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { api, getApiError } from "@/lib/api";
import { Event, User } from "@/lib/types";
import { ArrowLeft, CalendarDays, MapPin, Trash2, Users } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import Link from "next/link";
import { useParams, useRouter } from "next/navigation";
import { toast } from "sonner";
import useSWR from "swr";

const fetcher = (url: string) => api.get(url).then((r) => r.data);

export default function EventDetailPage() {
  const { id } = useParams<{ id: string }>();
  const router = useRouter();

  // Fetch event details
  // SWR = Stale-While-Revalidate
  // useSWR automatically refetches when:
  // - User switches browser tabs and comes back
  // - User reconnects to internet
  // - Component refocuses
  // - Every X seconds (configurable)
<<<<<<< HEAD
  const {
    data: event,
    error,
    mutate,
  } = useSWR<Event>(`/events/${id}`, fetcher);
=======
  const { data: event, error } = useSWR<Event>(`/events/${id}`, fetcher);
>>>>>>> b2b83c2 (Added add-attendee page, menus for profile and settings)
  const {
    data: attendees,
    error: attendeesError,
    isLoading: attendeesLoading,
    mutate: mutateAtt,
  } = useSWR<User[]>(`/events/${id}/attendees`, fetcher);
  if (error) toast.error(getApiError(error));

  const formatEventDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString("en-US", {
      weekday: "long",
      year: "numeric",
      month: "long",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  const getEventStatus = (dateString: string) => {
    const eventDate = new Date(dateString);
    const now = new Date();
    const diffHours = (eventDate.getTime() - now.getTime()) / (1000 * 60 * 80);

    if (diffHours < -2)
      return {
        status: "Past",
        color: "secondary",
      };
    if (diffHours < 2)
      return {
        status: "Live Now",
        color: "destructive",
      };
    if (diffHours < 24)
      return {
        status: "Starting Soon",
        color: "default",
      };
    return { status: "Upcoming", color: "outline" };
  };

  async function remove() {
    if (!window.confirm("Are you sure you want to delete this event?")) return;

    try {
      await api.delete(`/events/${id}`);
      toast.success("Event deleted successfully.");
      router.push("/events");
    } catch (e) {
      toast.error(getApiError(e));
    }
  }

  async function removeAttendee(uid: number) {
    if (!window.confirm("Remove this attendee from the event?")) return;
    try {
      await api.delete(`/events/${id}/attendees/${uid}`);
      toast.success("Attendee removed successfully.");
      mutateAtt(); // Refresh attendee list
    } catch (e) {
      toast.error(getApiError(e));
    }
  }

  if (!event) {
    return (
      <div className="flex items-center justify-center py-12">
        <div className="flex items-center gap-2 text-muted-foreground">
          <div className="h-4 w-4 animate-spin rounded-full border-2 border-primary border-t-transparent" />
          <span>Loading event details...</span>
        </div>
      </div>
    );
  }

  const eventStatus = getEventStatus(event.date);

  return (
    <div className="container mx-auto py-6 space-y-6">
      {/* Header */}
      <div className="flex items-center gap-4">
        <Button asChild variant="ghost" size="sm">
          <Link href="/events">
            <ArrowLeft />
            Back to Events
          </Link>
        </Button>
      </div>

      {/* Event Details Card */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between gap-4">
            <div className="space-y-2">
              <CardTitle className="text-2xl">{event.name}</CardTitle>
              <Badge
                variant={
                  eventStatus.color as
                    | "secondary"
                    | "destructive"
                    | "default"
                    | "outline"
                }
              >
                {eventStatus.status}
              </Badge>
            </div>

            <div className="flex items-center gap-2">
              <Button
                asChild
                variant="outline"
                className="hover:bg-gray-50 transition-colors"
              >
                <Link href={`/events/${id}/edit`}>Edit Event</Link>
              </Button>
              <Button
                variant="destructive"
                onClick={remove}
                className="hover:bg-red-600 transition-all duration-200 cursor-pointer shadow-md hover:shadow-lg hover:border-red-400 group relative overflow-hidden"
              >
                <Trash2 className="w-4 h-4 mr-2 group-hover:animate-pulse" />
                Delete Event
                {/* Optional: Add a subtle glow effect */}
                <div className="absolute inset-0 bg-red-400 opacity-0 group-hover:opacity-20 transition-opacity duration-200 rounded" />
              </Button>
            </div>
          </div>
        </CardHeader>

        <CardContent className="space-y-6">
          <p className="text-muted-foreground text-lg">{event.description}</p>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="flex items-center gap-3">
              <CalendarDays className="h-5 w-5 text-primary" />
              <div>
                <p className="font-medium">Date & Time</p>
                <p className="text-sm text-muted-foreground">
                  {formatEventDate(event.date)}
                </p>
              </div>
            </div>

            <div className="flex items-center gap-3">
              <MapPin className="h-5 w-5 text-primary" />
              <div>
                <p className="font-medium">Location</p>
                <p className="text-sm text-muted-foreground">
                  {event.location}
                </p>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Attendees Section */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <Users className="h-5 w-5" />
              <CardTitle>Attendees ({attendees?.length || 0})</CardTitle>
            </div>
            <Button asChild size="sm">
              <Link href={`/events/${id}/add-attendee`}>Add Attendee</Link>
            </Button>
          </div>
        </CardHeader>

        <CardContent>
          {attendeesLoading ? (
            <div className="flex items-center gap-2 text-muted-foreground">
              <div className="h-4 w-4 animate-spin rounded-full border-2 border-primary border-t-transparent" />
              <span>Loading attendees...</span>
            </div>
          ) : attendeesError ? (
            <div className="text-center py-8 text-red-500">
              <p>Failed to load attendees</p>
              <Button
                onClick={() => mutateAtt()}
                size="sm"
                variant="outline"
                className="mt-2"
              >
                Retry
              </Button>
            </div>
          ) : !attendees || attendees.length === 0 ? (
            <div className="text-center py-8">
              <Users className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
              <p className="text-muted-foreground">No attendees yet</p>
              <Button asChild className="mt-4" size="sm">
                <Link href={`/events/${id}/add-attendee`}>
                  Add First Attendee
                </Link>
              </Button>
            </div>
          ) : (
            // ✅ Show when we have attendees
            <div className="space-y-2">
              {attendees.map((attendee) => (
                <div
                  key={attendee.id}
                  className="flex items-center justify-between p-3 border rounded-lg"
                >
                  <div>
                    <p className="font-medium">{attendee.name}</p>
                    <p className="text-sm text-muted-foreground">
                      {attendee.email}
                    </p>
                  </div>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => removeAttendee(attendee.id)}
                    className="text-destructive hover:text-destructive"
                  >
                    Remove
                  </Button>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
