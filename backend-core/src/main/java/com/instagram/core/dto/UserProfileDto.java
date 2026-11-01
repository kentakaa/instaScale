package com.instagram.core.dto;

import com.instagram.core.model.User;

public record UserProfileDto(
    String id,
    String username,
    String email,
    String fullName,
    String bio,
    String avatarUrl,
    boolean isPrivate,
    boolean isVerified,
    long followerCount,
    long followingCount,
    long postCount
) {
    public static UserProfileDto from(User user) {
        return new UserProfileDto(
            user.getId(),
            user.getUsername(),
            user.getEmail(),
            user.getFullName(),
            user.getBio(),
            user.getAvatarUrl(),
            user.isPrivate(),
            user.isVerified(),
            user.getFollowerCount(),
            user.getFollowingCount(),
            user.getPostCount()
        );
    }
}
