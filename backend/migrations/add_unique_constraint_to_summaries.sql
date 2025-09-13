-- Migration to add UNIQUE constraint to summaries.note_id
-- This fixes the "Failed to save summary" error in production

-- First, remove any duplicate entries (keep the latest one)
DELETE FROM summaries 
WHERE id NOT IN (
    SELECT MAX(id) 
    FROM summaries 
    GROUP BY note_id
);

-- Add the UNIQUE constraint
ALTER TABLE summaries ADD CONSTRAINT summaries_note_id_unique UNIQUE (note_id);

-- Add an index for better performance
CREATE INDEX IF NOT EXISTS idx_summaries_note_id_unique ON summaries(note_id);
