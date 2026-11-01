package com.instagram.core.service;

import com.instagram.core.dto.AuthResponse;
import com.instagram.core.dto.LoginRequest;
import com.instagram.core.dto.RegisterRequest;
import com.instagram.core.dto.UserProfileDto;
import com.instagram.core.model.User;
import com.instagram.core.repository.UserRepository;
import com.instagram.core.security.JwtService;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;

import java.util.HashMap;
import java.util.Map;

@Service
public class AuthService {

    private final UserRepository userRepository;
    private final PasswordEncoder passwordEncoder;
    private final JwtService jwtService;
    private final AuthenticationManager authenticationManager;

    public AuthService(
            UserRepository userRepository,
            PasswordEncoder passwordEncoder,
            JwtService jwtService,
            AuthenticationManager authenticationManager) {
        this.userRepository = userRepository;
        this.passwordEncoder = passwordEncoder;
        this.jwtService = jwtService;
        this.authenticationManager = authenticationManager;
    }

    public AuthResponse register(RegisterRequest request) {
        String normalizedUsername = request.username().trim().toLowerCase();
        String contact = request.contact().trim().toLowerCase();

        if (userRepository.existsByUsername(normalizedUsername)) {
            throw new IllegalArgumentException("Username '" + normalizedUsername + "' is already taken.");
        }

        if (userRepository.existsByEmail(contact)) {
            throw new IllegalArgumentException("An account with email/contact '" + contact + "' already exists.");
        }

        User user = new User();
        user.setUsername(normalizedUsername);
        user.setEmail(contact);
        user.setFullName(request.fullName().trim());
        user.setPasswordHash(passwordEncoder.encode(request.password()));

        User savedUser = userRepository.save(user);

        Map<String, Object> claims = new HashMap<>();
        claims.put("userId", savedUser.getId());
        claims.put("email", savedUser.getEmail());

        String accessToken = jwtService.generateAccessToken(savedUser.getUsername(), claims);
        String refreshToken = jwtService.generateRefreshToken(savedUser.getUsername());

        return AuthResponse.of(
                accessToken,
                refreshToken,
                savedUser.getId(),
                savedUser.getUsername(),
                savedUser.getEmail(),
                savedUser.getFullName()
        );
    }

    public AuthResponse login(LoginRequest request) {
        String identifier = request.usernameOrEmail().trim().toLowerCase();

        Authentication authentication = authenticationManager.authenticate(
                new UsernamePasswordAuthenticationToken(identifier, request.password())
        );

        String principalUsername = authentication.getName();
        User user = userRepository.findByUsername(principalUsername)
                .or(() -> userRepository.findByEmail(principalUsername))
                .orElseThrow(() -> new IllegalArgumentException("User not found"));

        Map<String, Object> claims = new HashMap<>();
        claims.put("userId", user.getId());
        claims.put("email", user.getEmail());

        String accessToken = jwtService.generateAccessToken(user.getUsername(), claims);
        String refreshToken = jwtService.generateRefreshToken(user.getUsername());

        return AuthResponse.of(
                accessToken,
                refreshToken,
                user.getId(),
                user.getUsername(),
                user.getEmail(),
                user.getFullName()
        );
    }

    public UserProfileDto getCurrentUser(String username) {
        User user = userRepository.findByUsername(username)
                .orElseThrow(() -> new IllegalArgumentException("User not found with username: " + username));
        return UserProfileDto.from(user);
    }
}
