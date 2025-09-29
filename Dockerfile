# --- Build Stage ---
# Use the official Go image as a builder.
# Using '-alpine' makes the image smaller.
FROM golang:1.25-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy dependency files and download dependencies.
# This is done in a separate step to leverage Docker's layer caching.
# Dependencies are only re-downloaded if go.mod or go.sum changes.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application.
# CGO_ENABLED=0 creates a statically linked binary.
# GOOS=linux ensures it's built for the Linux OS inside the container.
RUN CGO_ENABLED=0 GOOS=linux go build -o /horus-app .

# --- Final Stage ---
# Use a minimal, empty base image for the final container.
# 'scratch' is the smallest possible image.
FROM scratch AS final

# Set the working directory
WORKDIR /

# Copy only the compiled binary from the 'builder' stage.
# This makes our final image incredibly small and secure.
COPY --from=builder /horus-app /horus-app

# Expose port 8080 to the outside world.
EXPOSE 8080

# The command to run when the container starts.
ENTRYPOINT ["/horus-app"]