CREATE TABLE IF NOT EXISTS admission_periods (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    school_year VARCHAR(50) NOT NULL,
    start_date DATE NULL,
    end_date DATE NULL,
    exam_fee DECIMAL(15,2) DEFAULT 0,
    is_open BOOLEAN DEFAULT FALSE,
    notes TEXT,
    created_at DATETIME(3) NULL,
    updated_at DATETIME(3) NULL
);
