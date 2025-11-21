"use client";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useAuth } from "@/contexts/AuthContext";
import { api, getApiError } from "@/lib/api";
import { zodResolver } from "@hookform/resolvers/zod";
import { ArrowLeft, Users } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import React, { use, useState } from "react";
import { useForm, type Resolver } from "react-hook-form";
import { toast } from "sonner";
import useSWR from "swr";
import z from "zod";

type EventSummary = {
  id: number;
  name: string;
  location: string;
  date: string;
};

const schema = z.object({
  userId: z.coerce.number().int().positive("Please provide a valid user ID"),
});

const fetcher = (url: string) => api.get(url).then((r) => r.data);

const AddAtendeePage = ({ params }: { params: Promise<{ id: string }> }) => {
  const resolvedParams = use(params);
  const eventId = Number(resolvedParams.id);
  const router = useRouter();
  const { isAuthed } = useAuth();
  const { data: event, isLoading } = useSWR<EventSummary>(
    isAuthed ? `/events/${eventId}` : null,
    fetcher
  );

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema) as Resolver<z.infer<typeof schema>>,
    defaultValues: { userId: undefined },
  });

  const onSubmit = async (values: z.infer<typeof schema>) => {
    try {
      await api.post(`/events/${eventId}/attendees/${values.userId}`);
      toast.success("Attendee added successfully");
      router.push(`/events/${eventId}`);
    } catch (error) {
      toast.error(getApiError(error));
    }
  };

  if (!isAuthed) {
    return (
      <div className="container mx-auto py-10 max-w-2xl">
        <Card>
          <CardHeader>
            <CardTitle>Authentication required</CardTitle>
            <CardDescription>Log in to manage attendees.</CardDescription>
          </CardHeader>
          <CardContent>
            <Button asChild>
              <Link href="/login">Go to login</Link>
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="container mx-auto py-8 max-w-2xl space-y-6">
      <div className="flex items-center gap-3">
        <Button asChild variant={"ghost"} size={"sm"}>
          <Link href={`/events/${eventId}`}>
            <ArrowLeft className="h-4 w-4" />{" "}
            <span className="sr-only">Back</span>
          </Link>
        </Button>
        <h1 className="text-2xl font-bold">Add attendee</h1>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>
            {isLoading ? "Loading event..." : event?.name ?? "Event not found"}
          </CardTitle>
          {event && (
            <CardDescription>
              <p>{event.location}</p>
              <p className="text-sm text-muted-foreground">
                {new Date(event.date).toLocaleString()}
              </p>
            </CardDescription>
          )}
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="rounded-lg border bg-muted/40 p-4 text-sm flex items-center gap-2">
            <Users className="h-4 w-4 text-muted-foreground" />
            Enter the numeric ID of the user you want to invite.
          </div>

          <form className="space-y-4" onSubmit={form.handleSubmit(onSubmit)}>
            <div className="space-y-2">
              <Label htmlFor="userId">User ID</Label>
              <Input
                id="userId"
                type="number"
                placeholder="e.g. 42"
                {...form.register("userId")}
              />
              {form.formState.errors.userId && (
                <p className="text-sm text-destructive">
                  {form.formState.errors.userId.message}
                </p>
              )}
            </div>
            <Button
              type="submit"
              className="w-full"
              disabled={form.formState.isSubmitting}
            >
              {form.formState.isSubmitting ? "Adding..." : "Add attendee"}
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  );
};

export default AddAtendeePage;
