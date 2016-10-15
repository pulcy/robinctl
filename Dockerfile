FROM alpine:3.4

# Add files
ADD ./robinctl /app/

# Start the app
ENTRYPOINT ["/app/robinctl"]
