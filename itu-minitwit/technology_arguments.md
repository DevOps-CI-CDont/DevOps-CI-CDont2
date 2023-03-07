# Minitwit 2.0

## Frontend choice
We choose Next.js as our frontend service, due to it's many build in features, such as:
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
It is a new language for all group members, so we can learn it together. We thought it would be an interesting language to build an API in. 
We first expected to use Go with Gorilla Mux, but we found out that that library has been archived.

## Database choice
_to be described_ 

# DigitalOcean  
Droplet setup .. 
Container registry ..

## CI/CD choice
We chose to use GitHub Actions for our CI/CD pipeline, because it is free and is what the course says to use.
Pytest .. 
go test ..
Watchtower ..
