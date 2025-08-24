"use client";
import { api } from "@/lib/api";
import useSwr from "swr";

const fetcher = (url: string) => api.get(url).then((r) => r.data);

export default function EventsPage() {}
