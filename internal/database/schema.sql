BEGIN;

-- Таблица пользователей
CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       login VARCHAR(255) UNIQUE NOT NULL,
                       password_hash VARCHAR(255) NOT NULL, --Будем хранить хэш, а не сам пароль
                       email VARCHAR(255) UNIQUE,
                       name VARCHAR(255),
                       surname VARCHAR(255),
                       registered_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Таблица авторизаций
CREATE TABLE auth_logs (
                           id SERIAL PRIMARY KEY,
                           user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
                           is_authorized BOOLEAN DEFAULT FALSE,
                           authorized_at TIMESTAMP WITHOUT TIME ZONE,
                           logged_out_at TIMESTAMP WITHOUT TIME ZONE
);

-- Таблица доступных вариантов тестов
CREATE TABLE test_variants (
                               id SERIAL PRIMARY KEY,
                               name VARCHAR(255) NOT NULL
);

-- Таблица заданий для вариантов
CREATE TABLE test_questions (
                                id SERIAL PRIMARY KEY,
                                variant_id INTEGER REFERENCES test_variants(id) ON DELETE SET NULL,
                                question_number INTEGER,
                                question_text TEXT NOT NULL,
                                correct_answer TEXT NOT NULL,
                                option_1 TEXT NOT NULL,
                                option_2 TEXT NOT NULL,
                                option_3 TEXT NOT NULL,
                                option_4 TEXT NOT NULL,
                                question_type VARCHAR(255) DEFAULT 'single_choice'
);

-- Таблица начала тестирования
CREATE TABLE user_tests (
                            id SERIAL PRIMARY KEY,
                            user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
                            variant_id INTEGER REFERENCES test_variants(id) ON DELETE SET NULL,
                            started_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                            total_questions INTEGER
);

-- Таблица ответов пользователей на задания
CREATE TABLE user_answers (
                              id SERIAL PRIMARY KEY,
                              user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
                              test_id INTEGER REFERENCES user_tests(id) ON DELETE SET NULL,
                              question_id INTEGER REFERENCES test_questions(id) ON DELETE SET NULL,
                              answer TEXT,
                              is_correct BOOLEAN
);

-- Таблица результатов тестирования
CREATE TABLE test_results (
                              id SERIAL PRIMARY KEY,
                              test_id INTEGER REFERENCES user_tests(id) ON DELETE SET NULL,
                              correct_answers_count INTEGER
);

-- Таблица для журналирования
CREATE TYPE action_type AS ENUM ('login', 'start_test', 'end_test', 'other');

CREATE TABLE user_logs (
                           id SERIAL PRIMARY KEY,
                           user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
                           action action_type,
                           timestamp TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                           details TEXT
);

-- Таблица сессий
CREATE TABLE sessions (
                          id SERIAL PRIMARY KEY,
                          token VARCHAR(256) NOT NULL UNIQUE,
                          user_id INT NOT NULL,
                          created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


-- Индексы
CREATE INDEX idx_user_login ON users(login);
CREATE INDEX idx_user_test ON user_tests(user_id, variant_id);
CREATE INDEX idx_question_variant ON test_questions(variant_id);

COMMIT;
