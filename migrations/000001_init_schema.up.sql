-- migrations/000001_init_schema.up.sql
-- Enrollment System Final Database Schema

-- Provinces
CREATE TABLE IF NOT EXISTS provinces (
    id INT AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- WardUnit
CREATE TABLE IF NOT EXISTS ward_units (
    id INT AUTO_INCREMENT PRIMARY KEY,
    province_id INT NOT NULL,
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (province_id) REFERENCES provinces(id) ON DELETE CASCADE
);

-- UserAccount
CREATE TABLE IF NOT EXISTS user_accounts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(100) UNIQUE,
    phone_number VARCHAR(20),
    full_name VARCHAR(100) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- AdminUser
CREATE TABLE IF NOT EXISTS admin_users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    status ENUM('active', 'inactive') DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- RoleGroup
CREATE TABLE IF NOT EXISTS role_groups (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- AdminUserRoleGroup
CREATE TABLE IF NOT EXISTS admin_user_role_groups (
    admin_user_id INT NOT NULL,
    role_group_id INT NOT NULL,
    PRIMARY KEY (admin_user_id, role_group_id),
    FOREIGN KEY (admin_user_id) REFERENCES admin_users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_group_id) REFERENCES role_groups(id) ON DELETE CASCADE
);

-- RoleGroupPermission
CREATE TABLE IF NOT EXISTS role_group_permissions (
    role_group_id INT NOT NULL,
    permission_code VARCHAR(100) NOT NULL,
    PRIMARY KEY (role_group_id, permission_code),
    FOREIGN KEY (role_group_id) REFERENCES role_groups(id) ON DELETE CASCADE
);

-- Menu
CREATE TABLE IF NOT EXISTS menus (
    id INT AUTO_INCREMENT PRIMARY KEY,
    parent_id INT DEFAULT NULL,
    name VARCHAR(100) NOT NULL,
    path VARCHAR(255) NOT NULL,
    icon VARCHAR(100),
    sort_order INT DEFAULT 0,
    FOREIGN KEY (parent_id) REFERENCES menus(id) ON DELETE SET NULL
);

-- AdmissionPeriod
CREATE TABLE IF NOT EXISTS admission_periods (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    status ENUM('upcoming', 'active', 'closed') DEFAULT 'upcoming',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Application
CREATE TABLE IF NOT EXISTS applications (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_account_id INT NOT NULL,
    admission_period_id INT NOT NULL,
    tracking_code VARCHAR(50) NOT NULL UNIQUE,
    full_name VARCHAR(100) NOT NULL,
    dob DATE NOT NULL,
    gender ENUM('male', 'female', 'other') NOT NULL,
    address TEXT NOT NULL,
    ward_unit_id INT NOT NULL,
    status ENUM('draft', 'submitted', 'under_review', 'approved', 'rejected') DEFAULT 'draft',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_account_id) REFERENCES user_accounts(id),
    FOREIGN KEY (admission_period_id) REFERENCES admission_periods(id),
    FOREIGN KEY (ward_unit_id) REFERENCES ward_units(id)
);

-- AcademicRecord
CREATE TABLE IF NOT EXISTS academic_records (
    id INT AUTO_INCREMENT PRIMARY KEY,
    application_id INT NOT NULL UNIQUE,
    school_name VARCHAR(255) NOT NULL,
    graduation_year INT NOT NULL,
    gpa DECIMAL(4,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (application_id) REFERENCES applications(id) ON DELETE CASCADE
);

-- ApplicationDocument
CREATE TABLE IF NOT EXISTS application_documents (
    id INT AUTO_INCREMENT PRIMARY KEY,
    application_id INT NOT NULL,
    document_type VARCHAR(50) NOT NULL,
    file_url VARCHAR(500) NOT NULL,
    status ENUM('pending', 'verified', 'rejected') DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (application_id) REFERENCES applications(id) ON DELETE CASCADE
);

-- ExamRoom
CREATE TABLE IF NOT EXISTS exam_rooms (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    capacity INT NOT NULL,
    location VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Examiner
CREATE TABLE IF NOT EXISTS examiners (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    phone VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- ExamRoomAssignment
CREATE TABLE IF NOT EXISTS exam_room_assignments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    application_id INT NOT NULL,
    exam_room_id INT NOT NULL,
    seat_number VARCHAR(20),
    FOREIGN KEY (application_id) REFERENCES applications(id) ON DELETE CASCADE,
    FOREIGN KEY (exam_room_id) REFERENCES exam_rooms(id) ON DELETE CASCADE
);

-- ExaminerAssignment
CREATE TABLE IF NOT EXISTS examiner_assignments (
    examiner_id INT NOT NULL,
    exam_room_id INT NOT NULL,
    role VARCHAR(50) DEFAULT 'proctor',
    PRIMARY KEY (examiner_id, exam_room_id),
    FOREIGN KEY (examiner_id) REFERENCES examiners(id) ON DELETE CASCADE,
    FOREIGN KEY (exam_room_id) REFERENCES exam_rooms(id) ON DELETE CASCADE
);

-- Subject
CREATE TABLE IF NOT EXISTS subjects (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(20) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- AdmissionPeriodSubject
CREATE TABLE IF NOT EXISTS admission_period_subjects (
    admission_period_id INT NOT NULL,
    subject_id INT NOT NULL,
    coefficient DECIMAL(3,1) DEFAULT 1.0,
    PRIMARY KEY (admission_period_id, subject_id),
    FOREIGN KEY (admission_period_id) REFERENCES admission_periods(id) ON DELETE CASCADE,
    FOREIGN KEY (subject_id) REFERENCES subjects(id) ON DELETE CASCADE
);

-- ApplicationExamScore
CREATE TABLE IF NOT EXISTS application_exam_scores (
    id INT AUTO_INCREMENT PRIMARY KEY,
    application_id INT NOT NULL,
    subject_id INT NOT NULL,
    score DECIMAL(4,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (application_id) REFERENCES applications(id) ON DELETE CASCADE,
    FOREIGN KEY (subject_id) REFERENCES subjects(id) ON DELETE CASCADE,
    UNIQUE (application_id, subject_id)
);

-- ApplicationResult
CREATE TABLE IF NOT EXISTS application_results (
    id INT AUTO_INCREMENT PRIMARY KEY,
    application_id INT NOT NULL UNIQUE,
    total_score DECIMAL(6,2),
    is_admitted BOOLEAN DEFAULT FALSE,
    note TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (application_id) REFERENCES applications(id) ON DELETE CASCADE
);
