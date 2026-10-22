CREATE TABLE IF NOT EXISTS jobs (
    id TEXT PRIMARY KEY,
    company TEXT NOT NULL,
    job_title TEXT NOT NULL,
    job_url TEXT,
    date_applied TEXT,
    deadline TEXT,
    status TEXT NOT NULL DEFAULT 'Saved'
        CHECK (status IN (
            'Saved',
            'Applied',
            'Interview',
            'Offer',
            'Rejected'
        )),
    notes TEXT,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);