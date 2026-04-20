CREATE TABLE IF NOT EXISTS academic_records (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    application_id BIGINT UNSIGNED NOT NULL,
    grade_level INT NOT NULL,
    school_name VARCHAR(255) NOT NULL,
    school_location VARCHAR(255),
    school_year VARCHAR(50),
    academic_performance VARCHAR(100),
    conduct_rating VARCHAR(100),
    created_at DATETIME(3) NULL,
    updated_at DATETIME(3) NULL,
    
    UNIQUE INDEX idx_application_grade (application_id, grade_level),
    FOREIGN KEY (application_id) REFERENCES applications(id) ON DELETE CASCADE
);
