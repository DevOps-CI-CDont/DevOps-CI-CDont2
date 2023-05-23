export function Footer() {
	return (
		<div className="fixed w-screen bottom-0 border-t border-gray-500 dark:border-slate-900 shadow-md bg-white dark:bg-slate-900 ">
			<div className="flex flex-col justify-center items-center text-gray-600 dark:text-slate-100 h-20 text-sm">
				<span>Minitwit - A Go:Gin and Next.js application</span>
				<span>&copy; DevOps CI-CDon&apos;t boys</span>
			</div>
		</div>
	);
}
