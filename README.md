# Web Forum Project

## Overview
This project is a web-based forum that allows communication between users, categorization of posts, and interaction through likes and dislikes. It's designed to be a platform for sharing ideas, discussions, and content within a community.

## Features
- **User Authentication**: Register and login functionality with encrypted passwords.
- **Post & Comment System**: Registered users can create posts and comments.
- **Categories**: Posts can be associated with one or more categories.
- **Likes & Dislikes**: Registered users can like or dislike posts and comments.
- **Visibility**: All users (registered or not) can view posts and comments.
- **Filtering**: Users can filter posts by categories, their own posts, and liked posts (for registered users).
- **Docker Integration**: Packaged and isolated with Docker for easy deployment and scalability.

## Technologies
- **Front-End**: [Front-end technologies, e.g., HTML, CSS, JavaScript]
- **Back-End**: [Back-end technologies, e.g., Go]
- **Database**: SQLite
- **Authentication**: Cookie-based session management
- **Containerization**: Docker

## Setup Instructions
1. **Clone the Repository**

```bash
    git clone git@git.01.alem.school:styan/forum.git
    cd forum/
```


2. **Build and Run**
```bash
    make build
    make run
```

or

```bash
   go run ./cmd/
```


## Usage
- **Registration**: New users can register by providing an email, username, and password.
- **Login**: Users can log in with their email and password.
- **Creating Posts/Comments**: After logging in, users can create posts and associate them with categories. They can also comment on existing posts.
- **Interactions**: Users can like or dislike posts and comments.
- **Filtering**: Use the filter options to view posts by categories, your posts, or posts you've liked.

## Additional Notes
- Ensure Docker is installed and running on your machine before starting the project.
- All passwords are encrypted for security purposes.
- The session cookie has an expiration time set for enhanced security.


