# CloudBoard

### Description
A simple web app that synchronizes users clipboards across multiple devices. Users can seamlessly copy on one device and paste on another. Here's a demo of the application: https://youtu.be/8-8Q-dVBZpw

### Motivation
My main motivation to pursue this project was to learn **Go**, implement my own authentication system, and also familiarize myself with the Google Cloud Platform API's. I started the development using **Cloud SQL API**, **Cloud App Engine**, and **Cloud Memorystore for Redis**. I later switched to **Elephant SQL**, and **Redis Labs** due to credit limitations on GCP. The final application was also deployed on **Cloud Run**.

- I made my own JWT authentication system which enables users to signup, login, logout, and refresh their access tokens. Logging in creates a pair of refresh and access tokens which are sent using Set-Cookie headers. Refresh generates a new access token for the user and Logout deletes the tokens from the redis database.
- The SQL database was used to store user information and the redis database was used to store the tokens.
- Websockets were used to send clipboard data between the client and server. The server broadcasts the clipboard data to the client devices logged in with the same user. Redis pub/sub was used to achieve this.

I initially started coding the client side application in C++ but later switched to Go because of its available cross-platform library for clipboard access and ease of use.

### Installation
1. Create a database on [Elephant SQL](https://www.elephantsql.com/) and a redis database on [Redis Labs](https://app.redislabs.com/)
2. Create a .env file (Check sample.env for reference)
3. To run locally pull docker image and run the container
    ```
    docker pull ashutosh67/cloudboard-app
    ```
    ```
    docker run -p 8080:8080 --env-file=.env ashutosh67/cloudboard-app
    ```
4. To run on a server create a Google Cloud Platform project and enable Secret Manager API and Cloud Run API.
5. Run deploy-secrets.bash to deploy secrets to Secret Manager
6. Run deploy-run.bash to deploy the application to Cloud Run
7. Compile the client application and run it on your device (remember to change the server address in the client code)
    ```
    cd client && go build -o client
    ```
    ```
    ./client
    ```

### Resources
- https://codevoweb.com/how-to-properly-use-jwt-for-authentication-in-golang/
- https://www.youtube.com/watch?v=ma7rUS_vW9M
- https://cloud.google.com/sql/docs/mysql/connect-instance-app-engine
- https://github.com/GoogleCloudPlatform/golang-samples/blob/main/cloudsql/mysql/database-sql/connect_tcp.go
