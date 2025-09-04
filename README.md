# Chota URL Shortener 

Chota is a simple, full-stack custom URL shortener built with Go.
You can shorten long URLs into compact links, and optionally give them a custom name.
For example, turn https://www.google.com into urlchota.vercel.app/google.
I built this to learn Go, improve my backend skills, and create something actually useful.

## Why I Built This

> **This project was my hands-on journey to learn Go (Golang).**  
> I wanted to understand Go’s concurrency, web server capabilities, and how to build scalable backend services.  
> Everything here from the Redis integration to the custom handlers was written to deepen my understanding of Go.

## Features

-  Shorten any long URL to a compact, shareable link
-  Lets you choose a custom name for the shortened URL (like /myportfolio)
-  Fast redirects to original URLs
-  Persistent storage using Redis
-  Beautiful, responsive frontend (HTML/CSS/JS)
-  Secure, random short codes
-  Ready for cloud deployment (Vercel, etc.)

## Tech Stack

- **Backend:** Go (Golang)
- **Frontend:** HTML, CSS, JavaScript
- **Database:** Redis
- **Deployment:** Vercel (cloud-ready)

## How It Works

1. Enter your long URL in the form
2. (Optional) Choose a custom name (3–20 characters)
3. Click "Shorten URL"
4. Get a link like urlchota.vercel.app/yourname
5. Anyone who opens that link gets redirected

## Getting Started

```sh
# Clone the repo
git clone https://github.com/wqrzdn/ChotaURL-1_public
cd ChotaURL-1_public

# Install dependencies
go mod tidy

# Run locally
go run main.go
```

Open [http://localhost:3000](http://localhost:3000) in your browser.

## Demo

![Screenshot](screenshot.png) 

## Live Demo

 [View the live app here](https://urlchota.vercel.app/) 

## What I Learned

- Go’s web server and handler patterns
- Using Redis for fast, persistent storage
- Building RESTful APIs and custom middleware
- Frontend integration with backend APIs
- Cloud deployment best practices
