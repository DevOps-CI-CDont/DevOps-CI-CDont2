interface RouterProps {
	id: number;
	path: string;
	text: string;
}

export const router: RouterProps[] = [
	{
		id: 1,
		path: "/public",
		text: "Public timeline",
	},

	{
		id: 3,
		path: "/register",
		text: "Sign up",
	},
	{
		id: 4,
		path: "/login",
		text: "Sign in",
	},
];

export const authenticatedRouter: RouterProps[] = [
	{
		id: 1,
		path: "/public",
		text: "Public timeline",
	},
	{
		id: 2,
		path: "/",
		text: "My timeline",
	},
];
