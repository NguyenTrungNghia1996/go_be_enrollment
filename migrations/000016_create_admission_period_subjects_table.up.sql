CREATE TABLE IF NOT EXISTS admission_period_subjects (
    id INT AUTO_INCREMENT PRIMARY KEY,
    admission_period_id INT NOT NULL,
    subject_id INT NOT NULL,
    weight DECIMAL(5,2) NOT NULL DEFAULT 1.0,
    is_required BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (admission_period_id) REFERENCES admission_periods(id) ON DELETE CASCADE,
    FOREIGN KEY (subject_id) REFERENCES subjects(id) ON DELETE CASCADE,
    UNIQUE KEY idx_admission_period_subject (admission_period_id, subject_id)
);
