[![CodeFactor](https://www.codefactor.io/repository/github/notherealmarco/wasaphoto/badge?s=2a99529eb3b66797b3a0cae48a39232782ae6c1b)](https://www.codefactor.io/repository/github/notherealmarco/wasaphoto)

# WASAPhoto

*Keep in touch with your friends by sharing photos of special moments, thanks to WASAPhoto!*

*You can upload your photos directly from your PC, and they will be visible to everyone following you.*

(Live demo: [https://wasaphoto.marcorealacci.me](https://wasaphoto.marcorealacci.me))

---

```
This is my project for the Web And Software Architecture (WASA) class
```

### This project includes

* An API specification using the OpenAPI standard
* A backend written in the Go language
* A frontend in Vue.js
* Dockerfiles to deploy the backend and the frontend in containers.
  * Dockerfile.backend builds the container for the backend
  * Dockerfile.frontend builds the container for the frontend
  * Dockerfile.embedded builds the backend container, but the backend's webserver also delivers the frontend

### Before building

If you're building the project in production mode (see below), you need to specify the base URL for the backend in `vite.config.js`.


## Build & deploy

The only (officially) supported method is via Docker containers.

There are two supported methods.

#### Embedded build

This method is only recommended for testing purposes or instances with very few users (for performance reasons).

The following commands will build a single container to serve both frontend and backend.

```
docker build -t wasaphoto -f Dockerfile.embedded .
docker run -p 3000:3000 -v <path to data directory>:/data --name wasaphoto wasaphoto
```

Everything will be up and running on port 3000 (including the Web UI).


#### Production build

This method build two containers, one for the backend and a container that running nginx to serve the frontend.

This is very recommended on production envinoments.

1. Build and run the backend

   ```
   docker build -t wasabackend -f Dockerfile.backend .
   docker run -p 3000:3000 -v <path to data directory>:/data --name wasaphoto-backend wasabackend
   ```
2. Edit the `vite.config.js` file and replace `<your API URL>` with the backend's base URL.
3. Build and run the frontend

   ```
   docker build -t wasafrontend -f Dockerfile.frontend .
   docker run -p 8080:80 --name wasaphoto-frontend wasafrontend
   ```

The Web UI will be up and running on port 8080!

<your API URL>
