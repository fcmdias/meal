# Use an official Golang runtime as a parent image
FROM golang:latest

# Copy the current directory contents into the container at /app
COPY main .
# Copy the templates directory into the container
COPY templates ./templates


# Expose port 8080 for the container
EXPOSE 8080

# Run the Go app when the container starts
CMD ["./main"]
