-- Create projects table
CREATE TABLE IF NOT EXISTS projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create analysis_runs table
CREATE TABLE IF NOT EXISTS analysis_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    status TEXT NOT NULL DEFAULT 'Created',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    started_at TIMESTAMP,
    finished_at TIMESTAMP,
    input_type TEXT NOT NULL, -- 'zip' | 'snippet'
    input_meta JSONB,
    summary_json JSONB
);

-- Create issues table
CREATE TABLE IF NOT EXISTS issues (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    analysis_id UUID NOT NULL REFERENCES analysis_runs(id) ON DELETE CASCADE,
    severity TEXT NOT NULL,
    rule_code TEXT NOT NULL,
    message TEXT NOT NULL,
    file_path TEXT NOT NULL,
    line INT
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_issues_analysis_id ON issues(analysis_id);
CREATE INDEX IF NOT EXISTS idx_analysis_runs_project_id ON analysis_runs(project_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_analysis_runs_input_type ON analysis_runs(input_type);

