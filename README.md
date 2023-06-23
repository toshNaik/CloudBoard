# CloudBoard

### Description
A simple web app that synchronizes users clipboards across multiple devices. Users can seamlessly copy on one device and paste on another. Here's a demo of the application: https://youtu.be/8-8Q-dVBZpw

### Motivation
My main motivation to pursue this project was to learn **Go**, implement my own authentication system, and also familiarize myself with the Google Cloud Platform API's. I started the development using **Cloud SQL API**, **Cloud App Engine**, and **Cloud Memorystore for Redis**. I later switched to **Elephant SQL**, and **Redis Labs** due to credit limitations on GCP. The final application was also deployed on **Cloud Run**.

- I made my own JWT authentication system which enables users to signup, login, logout, and refresh their access tokens. Logging in creates a pair of refresh and access tokens which are sent using Set-Cookie headers. Refresh generates a new access token for the user and Logout deletes the tokens from the redis database.
- The SQL database was used to store user information and the redis database was used to store the tokens.
- Websockets were used to send clipboard data between the client and server. The server broadcasts the clipboard data to the client devices logged in with the same user. Redis pub/sub was used to achieve this.

I initially started coding the client side application in C++ but later switched to Go because of its available cross-platform library for clipboard access and ease of use.

### Resources
- https://codevoweb.com/how-to-properly-use-jwt-for-authentication-in-golang/
- https://www.youtube.com/watch?v=ma7rUS_vW9M
- https://cloud.google.com/sql/docs/mysql/connect-instance-app-engine
- https://github.com/GoogleCloudPlatform/golang-samples/blob/main/cloudsql/mysql/database-sql/connect_tcp.go
