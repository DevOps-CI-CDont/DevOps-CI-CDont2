# Minitwit 2.0 technology arguments

## Frontend choice

We chose Next.js as our frontend service, due to it's many build in features, such as:

- Code splitting
- Routing support
- Server Side Rendering
- Layouts and Components
- API routes
- Build in optimizations
- Support for middleware

Furthermore we have added Tailwindcss as a CSS utility class for easy styling.
Having Server Side Rendering (SSR) isn't necessary for the Minitwit application, but it allows us to fetch data and build a page on the server, which gives a much better score on Google Lighthouse and also removes Cumulative Layout Shift (CLS), which gives a better user experience.

Having code splitting, divides the js build files into smaller files, making the size smaller and therefore 'lighter' to fetch (Making the application faster).

Next.js handles routing and API routes in the pages directory (Next.js 13 = app directory). This removes the overhead of having to create a react router that points to the different page components. Having a dedicated API routes folder, allows us to create API routes, that can fetch our backend data with ease. The API routes folder removes many issues and privacy related data, such as enviornment variables.

The Minitwit application has a user authentication system. With Next.js we can create a middleware that removes unauthenticated users from accessing pages, they're not allowed to view.

## Backend choice

We chose Go with Gin to make a RESTful API because:
It is a new language for all group members, so we can learn it together. We thought it would be an interesting language to build an API in. It is a pretty modern language, that is said to be built for concurrency, it's fast.  
*"The story goes that Google engineers designed Go while waiting for other programs to compile. Their frustration at their toolset forced them to rethink system programming from the ground up, creating a lean, mean, and compiled solution that allows for massive multithreading, concurrency, and performance under pressure."* - <https://stackoverflow.blog/2020/11/02/go-golang-learn-fast-programming-languages/>
We first expected to use Go with Gorilla Mux, but we found out that that library has been archived - and so decided not to use Mux.

## Database choice

*to be described more* 

At first we "inherited a local" database setup (a db file inside backend/tmp), and we didn't change that when we started containerizing the application. This was a very bad setup for real production data, as the database would be lost when the container was restarted.

We have since changed the database to a PostgreSQL database on a DigitalOcean..

## Cloud provider - DigitalOcean

We went with DigitalOcean as a cloud provider since it has \$200 in free credits for students.
We created a Droplet setup manually, at first a very basic one with only 1 CPU and 1GB of RAM (7$ / month).
We have since upgraded the Droplet's hardware to 2GB RAM, after noticing that memory usage would steadily increase and eventually crash the backend of our minitwit (API for simulator and the app itself).
We chose to use Digital Ocean's own Container registry as we guessed that it might be best compatible with a Droplet on the same platform.

## CI/CD choice

We chose to use GitHub Actions for our CI/CD pipeline, because it is free and is what the course says to use. of course it integrates quite well with code already on Github.
In our CI workflow we run:

- Pytest for the simulator tests
- Go test for tests for application API

Our CD workflow runs if the CI workflow succeeds (avoid deploying broken code). In our CD pipeline we run:

- build images for backend and frontend
- tag images
- push images to DigitalOcean's container registry

On our droplet we run Watchtower (*"A container-based solution for automating Docker container base image updates."*) to deploy the latest images to the Droplet.
Watchtower is configured to pull the latest images from our container registry and restart the containers with the new images (should have extremely little downtime). Inspecting our image registry at a given interval.

## Domain

We have our domain purchased for free through Name.com, although we use DigitalOcean name servers. We utilize iptables in the linux droplet to redirect all HTTP traffic to port 3000 (where the next.js frontend is hosted).
 
## Static code analysis tools

TODO: write more

- zod
- eslint
- sonarqube
- husky (precommit hooks)

## Monitoring

TODO: write more

Prometheus + Grafana

## Logging setup

TODO: write more

EFK stack (as the exercises)
