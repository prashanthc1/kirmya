"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { profileClient, Profile, PublicProfileResponse } from "@/lib/api/profile";

export function useProfile() {
  return useQuery<Profile>({
    queryKey: ["profile", "me"],
    queryFn: () => profileClient.getMe(),
  });
}

export function usePublicProfile(userId: string) {
  return useQuery<PublicProfileResponse>({
    queryKey: ["profile", "public", userId],
    queryFn: () => profileClient.getPublicProfile(userId),
    enabled: !!userId,
  });
}

export function useUpdateProfileSection() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ section, data }: { section: string; data: any }) =>
      profileClient.patchSection(section, data),
    onSuccess: (updatedProfile) => {
      queryClient.invalidateQueries({ queryKey: ["profile", "me"] });
      if (updatedProfile && updatedProfile.user_id) {
        queryClient.invalidateQueries({
          queryKey: ["profile", "public", updatedProfile.user_id],
        });
      }
    },
  });
}
