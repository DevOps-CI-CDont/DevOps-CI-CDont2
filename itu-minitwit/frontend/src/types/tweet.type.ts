export type Tweet = {
	CreatedAt: string;
	DeledtedAt: string | null;
	UpdatedAt: string | null;
	ID: number;
	author_id: number;
	author_name: string;
	text: string;
	pub_date: number;
	flagged: number;
};
