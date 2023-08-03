FROM alpine:latest

# Creates an app directory to hold your app’s source code
WORKDIR /

# Copies everything from your root directory into /app
ADD ./assets ./assets
ADD ./view ./view
COPY ./game-server ./

# Tells Docker which network port your container listens on
EXPOSE 8080

# Specifies the executable command that runs when the container starts
CMD [ "/game-server" ]