package com.instagram.core.dto;

public record AuthResponse(
    String token,
    String refreshToken,
    String tokenType,
    String userId,
    String username,
    String email,
    String fullName
) {
    public static AuthResponse of(String token, String refreshToken, String userId, String username, String email, String fullName) {
        return new AuthResponse(token, refreshToken, "Bearer", userId, username, email, fullName);
    }
}
