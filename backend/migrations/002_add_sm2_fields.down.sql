-- ============================================================================
-- Rollback SM-2 Algorithm Fields from user_words
-- Migration 002 Down
-- ============================================================================

-- Drop indexes
DROP INDEX IF EXISTS idx_user_words_review_priority;
DROP INDEX IF EXISTS idx_user_words_next_review;

-- Remove SM-2 fields from user_words table
ALTER TABLE user_words
DROP COLUMN IF EXISTS next_review_at,
DROP COLUMN IF EXISTS last_reviewed_at,
DROP COLUMN IF EXISTS repetitions,
DROP COLUMN IF EXISTS interval,
DROP COLUMN IF EXISTS easiness_factor;
