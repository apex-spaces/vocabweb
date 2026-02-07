-- ============================================================================
-- VocabWeb Database Schema - Rollback Migration
-- PostgreSQL 15
-- Created: 2026-02-07
-- ============================================================================

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS schema_migrations;
DROP TABLE IF EXISTS study_plans;
DROP TABLE IF EXISTS exam_wordlists;
DROP TABLE IF EXISTS user_achievements;
DROP TABLE IF EXISTS achievements;
DROP TABLE IF EXISTS daily_stats;
DROP TABLE IF EXISTS review_logs;
DROP TABLE IF EXISTS user_word_tags;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS user_words;
DROP TABLE IF EXISTS groups;
DROP TABLE IF EXISTS words;
DROP TABLE IF EXISTS profiles;

-- Drop extension
DROP EXTENSION IF EXISTS "uuid-ossp";
