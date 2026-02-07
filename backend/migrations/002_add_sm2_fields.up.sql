-- ============================================================================
-- Add SM-2 Algorithm Fields to user_words
-- Migration 002
-- ============================================================================

-- Add SM-2 fields to user_words table
ALTER TABLE user_words
ADD COLUMN easiness_factor DECIMAL(4,2) DEFAULT 2.5 CHECK (easiness_factor >= 1.3),
ADD COLUMN interval INTEGER DEFAULT 0 CHECK (interval >= 0),
ADD COLUMN repetitions INTEGER DEFAULT 0 CHECK (repetitions >= 0),
ADD COLUMN last_reviewed_at TIMESTAMPTZ,
ADD COLUMN next_review_at TIMESTAMPTZ;

-- Create index for efficient due review queries
CREATE INDEX idx_user_words_next_review ON user_words(user_id, next_review_at) 
WHERE next_review_at IS NOT NULL;

-- Create index for forgetting probability calculation
CREATE INDEX idx_user_words_review_priority ON user_words(user_id, last_reviewed_at, easiness_factor, repetitions)
WHERE next_review_at IS NOT NULL;

COMMENT ON COLUMN user_words.easiness_factor IS 'SM-2 easiness factor (1.3-2.5+), higher = easier to remember';
COMMENT ON COLUMN user_words.interval IS 'Days until next review (SM-2 interval)';
COMMENT ON COLUMN user_words.repetitions IS 'Number of consecutive correct reviews';
COMMENT ON COLUMN user_words.last_reviewed_at IS 'Timestamp of last review';
COMMENT ON COLUMN user_words.next_review_at IS 'Calculated next review timestamp';
