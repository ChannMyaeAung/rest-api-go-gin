"use client";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { api, getApiError } from "@/lib/api";
import { eventSchema } from "@/lib/schema";
import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";
import Link from "next/link";
import { ArrowLeft } from "lucide-react";

export default function NewEventPage() {
  const router = useRouter();
  const form = useForm<z.infer<typeof eventSchema>>({
    resolver: zodResolver(eventSchema),
    defaultValues: {
      name: "",
      description: "",
      location: "",
      date: "",
    },
  });

  async function onSubmit(values: z.infer<typeof eventSchema>) {
    try {
      const eventDate = new Date(values.date).toISOString();

      const payload = {
        name: values.name,
        description: values.description,
        location: values.location,
        date: eventDate, // Send as ISO string
      };

      console.log("Payload being sent:", payload); // Debug log

      const { data } = await api.post(`/events`, payload);

      console.log("API response:", data); // Debug log

      toast.success("Event created successfully!");
      router.push(`/events/${data.event.id}`);
    } catch (e) {
      console.error("Error creating event:", e); // Debug log
      toast.error(getApiError(e));
    }
  }

  return (
    <div className="container mx-auto py-6 max-w-2xl">
      {/* Header */}
      <div className="flex items-center gap-4 mb-6">
        <Button asChild variant="ghost" size="sm">
          <Link href="/events">
            <ArrowLeft className="mr-2 h-4 w-4" />
          </Link>
        </Button>
        <h1 className="text-2xl font-bold">Create New Event</h1>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Event Details</CardTitle>
        </CardHeader>

        <CardContent>
          <form className="space-y-6" onSubmit={form.handleSubmit(onSubmit)}>
            <div className="space-y-2">
              <Label htmlFor="name">Event Name</Label>
              <Input
                id="name"
                placeholder="Enter event name"
                {...form.register("name")}
              />
              {/* ✅ Show validation errors */}
              {form.formState.errors.name && (
                <p className="text-sm text-red-500">
                  {form.formState.errors.name.message}
                </p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="description">Description</Label>
              <Textarea
                id="description"
                rows={4}
                placeholder="Describe your event"
                {...form.register("description")}
              />
              {form.formState.errors.description && (
                <p className="text-sm text-red-500">
                  {form.formState.errors.description.message}
                </p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="location">Location</Label>
              <Input
                id="location"
                placeholder="Where will this event take place?"
                {...form.register("location")}
              />
              {form.formState.errors.location && (
                <p className="text-sm text-red-500">
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
                // ✅ Set minimum date to current date
                min={new Date().toISOString().slice(0, 16)}
              />
              {form.formState.errors.date && (
                <p className="text-sm text-red-500">
                  {form.formState.errors.date.message}
                </p>
              )}
            </div>

            <Button
              type="submit"
              className="w-full"
              disabled={form.formState.isSubmitting}
            >
              {form.formState.isSubmitting
                ? "Creating Event..."
                : "Create Event"}
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
