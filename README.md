# url-shortener

My attempt at an URL Shortener

## What is a URL Shortener

A URL Shortener is a technique that takes a Uniform Resource Locator (URL) and shortens the length while still redirecting to the required page

## How does a URL Shortener work

- Every long URL is associated with a unique key, which is the part after the domain name ex.) https://tinyurl.com/a3b2st has a key of a3b2st
- Keys are case sensitive most of the time and using the wrong case may lead to a different destination URL

## Architecture

## Flow

A user will send a POST request to us and we should send back a shortURL back which will quickly redirect them to the same link
We should be able to check how many times that link is clicked
Each URL should have its own unique link and ID
When user clicks short URL it has to send to original URL right away

1. Check cache to see if it has seen short URL before
2. Check database if short URL is not in cache
3. Cache result for next time
4. System sends you to right URL
