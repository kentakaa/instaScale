package com.instagram.core.service;

import com.instagram.core.dto.UserProfileDto;
import com.instagram.core.model.User;
import com.instagram.core.repository.UserRepository;
import org.springframework.stereotype.Service;

@Service
public class UserService {

    private final UserRepository userRepository;

    public UserService(UserRepository userRepository) {
        this.userRepository = userRepository;
    }

    public UserProfileDto getProfile(String username) {
        User user = userRepository.findByUsername(username.toLowerCase())
                .orElseThrow(() -> new IllegalArgumentException("User not found: " + username));
        return UserProfileDto.from(user);
    }
}
