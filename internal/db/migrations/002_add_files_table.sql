-- Create analysis_files table to store file contents in database
CREATE TABLE IF NOT EXISTS analysis_files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    analysis_id UUID NOT NULL REFERENCES analysis_runs(id) ON DELETE CASCADE,
    file_path TEXT NOT NULL,
    content BYTEA NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create index for faster lookups
CREATE INDEX IF NOT EXISTS idx_analysis_files_analysis_id ON analysis_files(analysis_id);
CREATE INDEX IF NOT EXISTS idx_analysis_files_file_path ON analysis_files(analysis_id, file_path);

