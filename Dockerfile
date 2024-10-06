# Use the Arktix base image
FROM archlinux:latest

# Install Git and bash in the Arch environment
RUN pacman -Sy --noconfirm git bash pwd 

# Set the working directory
WORKDIR /action

# Copy all files into the working directory
COPY . .

# Make the Bash script executable
RUN chmod +x main.sh

# Run the Bash script
CMD ["bash", "/action/main.sh"]
