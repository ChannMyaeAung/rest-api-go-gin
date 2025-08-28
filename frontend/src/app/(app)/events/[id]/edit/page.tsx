"use client";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { api, getApiError } from "@/lib/api";
import { eventSchema } from "@/lib/schema";
import { Event } from "@/lib/types";
import { zodResolver } from "@hookform/resolvers/zod";
import Link from "next/link";
import { useParams, useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import useSWR from "swr";
import z from "zod";

const fetcher = (url: string) => api.get(url).then((r) => r.data);

// Convert ISO (UTC) to value usable by datetime-local (local time, no Z)
function isoToLocalInput(iso?: string): string {
  if (!iso) return "";
  const d = new Date(iso);

  // Adjust to local so that displayed time matches stored UTC time
  const off = d.getTimezoneOffset();
  const local = new Date(d.getTime() - off * 60000);
  return local.toISOString().slice(0, 16);
}

export default function EditEventPage() {
  const { id } = useParams<{ id: string }>();
  const router = useRouter();
  const {
    data: event,
    error,
    isLoading,
  } = useSWR<Event>(`/events/${id}`, fetcher);
  const [submitting, setSubmitting] = useState(false);

  const form = useForm<z.infer<typeof eventSchema>>({
    resolver: zodResolver(eventSchema),
    defaultValues: {
      name: "",
      description: "",
      location: "",
      date: "",
    },
  });

  // Update form values when data loads
  useEffect(() => {
    if (event) {
      form.reset({
        name: event.name,
        description: event.description,
        location: event.location,
        date: isoToLocalInput(event.date), // for datetime-local input
      });
    }
  }, [event, form]);

  async function onSubmit(values: z.infer<typeof eventSchema>) {
    try {
      setSubmitting(true);
      // datetime-local gives local time; convert to ISO UTC
      const isoDate = new Date(values.date).toISOString();
      const payload = {
        name: values.name,
        description: values.description,
        location: values.location,
        date: isoDate,
      };
      await api.put(`/events/${id}`, payload);
      toast.success("Event updated");
      router.push(`/events/${id}`);
    } catch (e) {
      toast.error(getApiError(e));
    } finally {
      setSubmitting(false);
    }
  }

  if (error) {
    return (
      <div>
        <p>Failed to load event: {getApiError(error)}</p>
        <Button variant="outline" onClick={() => location.reload()}>
          Retry
        </Button>
        <Button asChild>
          <Link href={`/events/${id}`}>Back</Link>
        </Button>
      </div>
    );
  }

  if (isLoading || !event)
    return (
      <div className="flex items-center justify-center py-16">
        <div className="flex items-center gap-2 text-muted-foreground">
          <div className="h-4 w-4 animate-spin rounded-full border-2 border-primary border-t-transparent" />
          <span>Loading the event...</span>
        </div>
      </div>
    );

  return (
    <div className="mx-auto max-w-xl py-6">
      <Card>
        <CardHeader>
          <CardTitle>Edit Event</CardTitle>
        </CardHeader>
        <CardContent>
          <form className="space-y-6" onSubmit={form.handleSubmit(onSubmit)}>
            <div className="space-y-2">
              <Label htmlFor="name">Name</Label>
              <Input id="name" {...form.register("name")} />
              {form.formState.errors.name && (
                <p className="text-xs text-red-500">
                  {form.formState.errors.name.message}
                </p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="description">Description</Label>
              <Textarea
                id="description"
                rows={5}
                {...form.register("description")}
              />
              {form.formState.errors.description && (
                <p className="text-xs text-red-500">
                  {form.formState.errors.description.message}
                </p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="location">Location</Label>
              <Input id="location" {...form.register("location")} />
              {form.formState.errors.location && (
                <p className="text-xs text-red-500">
                  {form.formState.errors.location.message}
                </p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="date">Date & Time</Label>
              <Input
                id="date"
                type="datetime-local"
                {...form.register("date")}
              />
              {form.formState.errors.date && (
                <p className="text-xs text-red-500">
                  {form.formState.errors.date.message}
                </p>
              )}
            </div>

            <div className="flex gap-3">
              <Button type="submit" disabled={submitting} className="min-w-32">
                {submitting ? "Saving..." : "Save Changes"}
              </Button>
              <Button
                type="button"
                variant="outline"
                onClick={() => router.push(`/events/${id}`)}
              >
                Cancel
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
