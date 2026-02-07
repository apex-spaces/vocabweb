-- ============================================================================
-- VocabWeb Database Schema - Initial Migration
-- PostgreSQL 15
-- Created: 2026-02-07
-- ============================================================================

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================================================
-- Table 1: profiles
-- User extended profile information
-- ============================================================================
CREATE TABLE profiles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id TEXT NOT NULL UNIQUE, -- From Google Identity Platform
    username VARCHAR(50) UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    avatar_url TEXT,
    daily_review_goal INTEGER DEFAULT 20 CHECK (daily_review_goal > 0),
    timezone VARCHAR(50) DEFAULT 'UTC',
    level INTEGER DEFAULT 1 CHECK (level >= 1),
    xp INTEGER DEFAULT 0 CHECK (xp >= 0),
    streak_days INTEGER DEFAULT 0 CHECK (streak_days >= 0),
    last_study_date DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_profiles_user_id ON profiles(user_id);
CREATE INDEX idx_profiles_username ON profiles(username);
CREATE INDEX idx_profiles_email ON profiles(email);

COMMENT ON TABLE profiles IS 'User extended profile information';
COMMENT ON COLUMN profiles.user_id IS 'Google Identity Platform user ID';
COMMENT ON COLUMN profiles.xp IS 'Experience points for gamification';
COMMENT ON COLUMN profiles.streak_days IS 'Consecutive study days';

-- ============================================================================
-- Table 2: words
-- Global word dictionary (shared across all users)
-- ============================================================================
CREATE TABLE words (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    word VARCHAR(100) NOT NULL UNIQUE,
    phonetic VARCHAR(100),
    definitions JSONB NOT NULL, -- Array of {pos: string, meaning: string, example: string}
    frequency_rank INTEGER, -- Lower = more common (1-100000)
    audio_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_words_word ON words(word);
CREATE INDEX idx_words_frequency_rank ON words(frequency_rank);
CREATE INDEX idx_words_definitions_gin ON words USING GIN(definitions);

COMMENT ON TABLE words IS 'Global word dictionary shared across all users';
COMMENT ON COLUMN words.definitions IS 'JSON array: [{pos: "noun", meaning: "...", example: "..."}]';
COMMENT ON COLUMN words.frequency_rank IS 'Word frequency rank (1=most common)';

-- ============================================================================
-- Table 3: groups
-- User-defined word groups/folders
-- ============================================================================
CREATE TABLE groups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id TEXT NOT NULL,
    name VARCHAR(100) NOT NULL,
    color VARCHAR(7) DEFAULT '#3B82F6', -- Hex color code
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_groups_user FOREIGN KEY (user_id) REFERENCES profiles(user_id) ON DELETE CASCADE,
    CONSTRAINT unique_user_group_name UNIQUE (user_id, name)
);

CREATE INDEX idx_groups_user_id ON groups(user_id);
CREATE INDEX idx_groups_sort_order ON groups(user_id, sort_order);

COMMENT ON TABLE groups IS 'User-defined word groups/folders';
COMMENT ON COLUMN groups.color IS 'Hex color code for UI display';

-- ============================================================================
-- Table 4: user_words
-- User's collected words with context
-- ============================================================================
CREATE TABLE user_words (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id TEXT NOT NULL,
    word_id UUID NOT NULL,
    group_id UUID,
    source_url TEXT,
    context_sentence TEXT,
    mastery_level INTEGER DEFAULT 0 CHECK (mastery_level BETWEEN 0 AND 5),
    is_mastered BOOLEAN DEFAULT FALSE,
    collected_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_user_words_user FOREIGN KEY (user_id) REFERENCES profiles(user_id) ON DELETE CASCADE,
    CONSTRAINT fk_user_words_word FOREIGN KEY (word_id) REFERENCES words(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_words_group FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE SET NULL,
    CONSTRAINT unique_user_word UNIQUE (user_id, word_id)
);

CREATE INDEX idx_user_words_user_id ON user_words(user_id);
CREATE INDEX idx_user_words_word_id ON user_words(word_id);
CREATE INDEX idx_user_words_group_id ON user_words(group_id);
CREATE INDEX idx_user_words_mastery ON user_words(user_id, mastery_level);
CREATE INDEX idx_user_words_mastered ON user_words(user_id, is_mastered);
CREATE INDEX idx_user_words_collected_at ON user_words(user_id, collected_at DESC);

COMMENT ON TABLE user_words IS 'User collected words with learning context';
COMMENT ON COLUMN user_words.mastery_level IS '0=new, 1-4=learning, 5=mastered';
COMMENT ON COLUMN user_words.context_sentence IS 'Original sentence where word was found';

-- ============================================================================
-- Table 5: tags
-- User-defined tags for organizing words
-- ============================================================================
CREATE TABLE tags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id TEXT NOT NULL,
    name VARCHAR(50) NOT NULL,
    color VARCHAR(7) DEFAULT '#10B981', -- Hex color code
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_tags_user FOREIGN KEY (user_id) REFERENCES profiles(user_id) ON DELETE CASCADE,
    CONSTRAINT unique_user_tag_name UNIQUE (user_id, name)
);

CREATE INDEX idx_tags_user_id ON tags(user_id);

COMMENT ON TABLE tags IS 'User-defined tags for organizing words';
COMMENT ON COLUMN tags.name IS 'Tag name (e.g., #GRE, #computer, #daily)';

-- ============================================================================
-- Table 6: user_word_tags
-- Many-to-many relationship between user_words and tags
-- ============================================================================
CREATE TABLE user_word_tags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_word_id UUID NOT NULL,
    tag_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_user_word_tags_user_word FOREIGN KEY (user_word_id) REFERENCES user_words(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_word_tags_tag FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE,
    CONSTRAINT unique_user_word_tag UNIQUE (user_word_id, tag_id)
);

CREATE INDEX idx_user_word_tags_user_word_id ON user_word_tags(user_word_id);
CREATE INDEX idx_user_word_tags_tag_id ON user_word_tags(tag_id);

COMMENT ON TABLE user_word_tags IS 'Many-to-many relationship between user words and tags';

-- ============================================================================
-- Table 7: review_logs
-- Spaced repetition review history with SM-2 algorithm fields
-- ============================================================================
CREATE TABLE review_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_word_id UUID NOT NULL,
    quality INTEGER NOT NULL CHECK (quality BETWEEN 0 AND 5),
    easiness_factor DECIMAL(4,2) DEFAULT 2.5 CHECK (easiness_factor >= 1.3),
    interval INTEGER DEFAULT 0 CHECK (interval >= 0),
    repetitions INTEGER DEFAULT 0 CHECK (repetitions >= 0),
    next_review_at TIMESTAMPTZ,
    reviewed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_review_logs_user_word FOREIGN KEY (user_word_id) REFERENCES user_words(id) ON DELETE CASCADE
);

CREATE INDEX idx_review_logs_user_word_id ON review_logs(user_word_id);
CREATE INDEX idx_review_logs_next_review ON review_logs(next_review_at);
CREATE INDEX idx_review_logs_reviewed_at ON review_logs(reviewed_at DESC);

COMMENT ON TABLE review_logs IS 'Spaced repetition review history with SM-2 algorithm';
COMMENT ON COLUMN review_logs.quality IS 'User rating: 0=complete blackout, 5=perfect recall';
COMMENT ON COLUMN review_logs.easiness_factor IS 'SM-2 easiness factor (1.3-2.5+)';
COMMENT ON COLUMN review_logs.interval IS 'Days until next review';
COMMENT ON COLUMN review_logs.repetitions IS 'Number of consecutive correct reviews';
COMMENT ON COLUMN review_logs.next_review_at IS 'Calculated next review timestamp';

-- ============================================================================
-- Table 8: daily_stats
-- Daily learning statistics snapshot
-- ============================================================================
CREATE TABLE daily_stats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id TEXT NOT NULL,
    stat_date DATE NOT NULL,
    new_words INTEGER DEFAULT 0 CHECK (new_words >= 0),
    reviewed INTEGER DEFAULT 0 CHECK (reviewed >= 0),
    mastered INTEGER DEFAULT 0 CHECK (mastered >= 0),
    streak_days INTEGER DEFAULT 0 CHECK (streak_days >= 0),
    xp_gained INTEGER DEFAULT 0 CHECK (xp_gained >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_daily_stats_user FOREIGN KEY (user_id) REFERENCES profiles(user_id) ON DELETE CASCADE,
    CONSTRAINT unique_user_stat_date UNIQUE (user_id, stat_date)
);

CREATE INDEX idx_daily_stats_user_id ON daily_stats(user_id);
CREATE INDEX idx_daily_stats_stat_date ON daily_stats(user_id, stat_date DESC);

COMMENT ON TABLE daily_stats IS 'Daily learning statistics snapshot for trend analysis';
COMMENT ON COLUMN daily_stats.new_words IS 'Number of new words collected today';
COMMENT ON COLUMN daily_stats.reviewed IS 'Number of words reviewed today';
COMMENT ON COLUMN daily_stats.mastered IS 'Number of words mastered today';

-- ============================================================================
-- Table 9: achievements
-- Achievement/badge definitions
-- ============================================================================
CREATE TABLE achievements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    icon VARCHAR(50) NOT NULL,
    condition_type VARCHAR(50) NOT NULL,
    condition_value INTEGER NOT NULL,
    xp_reward INTEGER DEFAULT 0 CHECK (xp_reward >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_achievements_name ON achievements(name);

COMMENT ON TABLE achievements IS 'Achievement/badge definitions';
COMMENT ON COLUMN achievements.condition_type IS 'e.g., streak_days, total_words, first_ocr';
COMMENT ON COLUMN achievements.condition_value IS 'Threshold value for unlocking';

-- ============================================================================
-- Table 10: user_achievements
-- User earned achievements
-- ============================================================================
CREATE TABLE user_achievements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id TEXT NOT NULL,
    achievement_id UUID NOT NULL,
    earned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_user_achievements_user FOREIGN KEY (user_id) REFERENCES profiles(user_id) ON DELETE CASCADE,
    CONSTRAINT fk_user_achievements_achievement FOREIGN KEY (achievement_id) REFERENCES achievements(id) ON DELETE CASCADE,
    CONSTRAINT unique_user_achievement UNIQUE (user_id, achievement_id)
);

CREATE INDEX idx_user_achievements_user_id ON user_achievements(user_id);
CREATE INDEX idx_user_achievements_achievement_id ON user_achievements(achievement_id);
CREATE INDEX idx_user_achievements_earned_at ON user_achievements(user_id, earned_at DESC);

COMMENT ON TABLE user_achievements IS 'User earned achievements/badges';

-- ============================================================================
-- Table 11: exam_wordlists
-- Exam-specific word lists (IELTS, TOEFL, GRE, etc.)
-- ============================================================================
CREATE TABLE exam_wordlists (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    exam_type VARCHAR(50) NOT NULL,
    word_id UUID NOT NULL,
    frequency_in_exam VARCHAR(20) DEFAULT 'medium',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_exam_wordlists_word FOREIGN KEY (word_id) REFERENCES words(id) ON DELETE CASCADE,
    CONSTRAINT unique_exam_word UNIQUE (exam_type, word_id),
    CONSTRAINT check_exam_type CHECK (exam_type IN ('ielts', 'toefl', 'gre', 'sat', 'cet4', 'cet6', 'kaoyan'))
);

CREATE INDEX idx_exam_wordlists_exam_type ON exam_wordlists(exam_type);
CREATE INDEX idx_exam_wordlists_word_id ON exam_wordlists(word_id);
CREATE INDEX idx_exam_wordlists_frequency ON exam_wordlists(exam_type, frequency_in_exam);

COMMENT ON TABLE exam_wordlists IS 'Exam-specific word lists with frequency data';
COMMENT ON COLUMN exam_wordlists.exam_type IS 'Exam type: ielts, toefl, gre, sat, cet4, cet6, kaoyan';
COMMENT ON COLUMN exam_wordlists.frequency_in_exam IS 'Frequency: high, medium, low';

-- ============================================================================
-- Table 12: study_plans
-- User exam preparation plans
-- ============================================================================
CREATE TABLE study_plans (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id TEXT NOT NULL,
    exam_type VARCHAR(50) NOT NULL,
    exam_date DATE NOT NULL,
    daily_target INTEGER DEFAULT 20 CHECK (daily_target > 0),
    total_words INTEGER DEFAULT 0 CHECK (total_words >= 0),
    completed_words INTEGER DEFAULT 0 CHECK (completed_words >= 0),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_study_plans_user FOREIGN KEY (user_id) REFERENCES profiles(user_id) ON DELETE CASCADE,
    CONSTRAINT check_study_plan_exam_type CHECK (exam_type IN ('ielts', 'toefl', 'gre', 'sat', 'cet4', 'cet6', 'kaoyan'))
);

CREATE INDEX idx_study_plans_user_id ON study_plans(user_id);
CREATE INDEX idx_study_plans_exam_type ON study_plans(exam_type);
CREATE INDEX idx_study_plans_exam_date ON study_plans(user_id, exam_date);
CREATE INDEX idx_study_plans_active ON study_plans(user_id, is_active);

COMMENT ON TABLE study_plans IS 'User exam preparation plans with daily targets';
COMMENT ON COLUMN study_plans.daily_target IS 'Number of words to learn per day';
COMMENT ON COLUMN study_plans.total_words IS 'Total words in the exam wordlist';
COMMENT ON COLUMN study_plans.completed_words IS 'Number of words already mastered';

-- ============================================================================
-- Migration tracking table
-- ============================================================================
CREATE TABLE IF NOT EXISTS schema_migrations (
    version VARCHAR(255) PRIMARY KEY,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE schema_migrations IS 'Tracks applied database migrations';
