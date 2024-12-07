# Build backend for Windows or Linux systems
FROM golang:1.23.4-bookworm AS backend-builder

# Run tests on backend based on the tests written
FROM debian:slim AS tester

# Build frontend for Windows or Linux systems
FROM node:slim AS frontend-builder
