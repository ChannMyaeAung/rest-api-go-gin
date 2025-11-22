"use client";
import { useAuth } from "@/contexts/AuthContext";
import { api, getApiError } from "@/lib/api";
import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import React, { useEffect } from "react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import z from "zod";
import { Button } from "@/components/ui/button";
import Link from "next/link";

const settingsSchema = z.object({
  name: z.string().min(2, "Name must be at least 2 characters long"),
  email: z.email("Invalid email address"),
  profile_picture: z
    .string()
    .url("Enter a valid URL")
    .optional()
    .or(z.literal("")),
});

type SettingsForm = z.infer<typeof settingsSchema>;

const SettingsPage = () => {
  const router = useRouter();
  const { isAuthed, isLoading, user, refreshUser, logout } = useAuth();

  const form = useForm<SettingsForm>({
    resolver: zodResolver(settingsSchema),
    defaultValues: {
      name: "",
      email: "",
      profile_picture: "",
    },
  });

  useEffect(() => {
    if (user) {
      form.reset({
        name: user.name ?? "",
        email: user.email ?? "",
        profile_picture: user.profile_picture ?? "",
      });
    }
  }, [user, form]);

  const handleSubmit = form.handleSubmit(async (values) => {
    try {
      await api.put("/auth/me", {
        name: values.name.trim(),
        email: values.email.trim(),
        profile_picture: values.profile_picture?.trim() || null,
      });
      toast.success("Profile updated successfully");
      await refreshUser();
    } catch (error) {
      toast.error("Failed to update profile. Please try again.");
    }
  });

  const handleDeleteAccount = async () => {
    if (
      !confirm(
        "Delete your account and all associated events? This action cannot be undone."
      )
    ) {
      return;
    }

    try {
      await api.delete("/auth/me");
      toast.success("Account deleted successfully");
      logout();
      router.replace("/register");
    } catch (error) {
      toast.error(getApiError(error));
    }
  };

  if (!isAuthed && !isLoading) {
    router.replace("/login");
    return null;
  }

  return (
    <div className="container mx-auto max-w-2xl space-y-6 py-10">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Account settings</h1>
          <p className="text-muted-foreground">
            Update your profile or delete your account.
          </p>
        </div>
        <Button variant="ghost" asChild>
          <Link href="/events">Back to events</Link>
        </Button>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Profile information</CardTitle>
          <CardDescription>
            Change your display name, email, or profile picture.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form className="space-y-4" onSubmit={handleSubmit}>
            <div className="space-y-2">
              <Label htmlFor="name">Name</Label>
              <Input id="name" {...form.register("name")} />
              {form.formState.errors.name && (
                <p className="text-sm text-destructive">
                  {form.formState.errors.name.message}
                </p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="email">Email</Label>
              <Input id="email" type="email" {...form.register("email")} />
              {form.formState.errors.email && (
                <p className="text-sm text-destructive">
                  {form.formState.errors.email.message}
                </p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="profile_picture">Profile picture URL</Label>
              <Input
                id="profile_picture"
                placeholder="https://example.com/avatar.jpg"
                {...form.register("profile_picture")}
              />
              {form.formState.errors.profile_picture && (
                <p className="text-sm text-destructive">
                  {form.formState.errors.profile_picture.message}
                </p>
              )}
            </div>

            <Button
              type="submit"
              className="w-full"
              disabled={form.formState.isSubmitting}
            >
              {form.formState.isSubmitting ? "Saving..." : "Save changes"}
            </Button>
          </form>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Danger zone</CardTitle>
          <CardDescription>
            Once you delete your account, there is no going back.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Button variant="destructive" onClick={handleDeleteAccount}>
            Delete account
          </Button>
        </CardContent>
      </Card>
    </div>
  );
};

export default SettingsPage;
