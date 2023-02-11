interface RouterProps {
  id: number;
  path?: string;
  onClick?: () => void;
  text: string;
  authenticated: boolean;
}

export const router: RouterProps[] = [
  {
    id: 1,
    path: "/public",
    text: "Public timeline",
    authenticated: false,
  },
  {
    id: 2,
    path: "/",
    text: "My timeline",
    authenticated: true,
  },
  {
    id: 3,
    path: "/register",
    text: "Sign up",
    authenticated: false,
  },
  {
    id: 4,
    path: "/login",
    text: "Sign in",
    authenticated: false,
  },
  {
    id: 5,

    text: "Sign out",
    authenticated: true,
  },
];
