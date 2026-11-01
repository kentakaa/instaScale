package com.instagram.core.controller;

import com.instagram.core.dto.ApiResponse;
import com.instagram.core.dto.UserProfileDto;
import com.instagram.core.service.UserService;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/api/v1/users")
public class UserController {

    private final UserService userService;

    public UserController(UserService userService) {
        this.userService = userService;
    }

    @GetMapping("/{username}")
    public ResponseEntity<ApiResponse<UserProfileDto>> getProfile(@PathVariable String username) {
        UserProfileDto profile = userService.getProfile(username);
        return ResponseEntity.ok(ApiResponse.success("Profile fetched successfully", profile));
    }
}
