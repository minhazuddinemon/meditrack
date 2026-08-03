-- 1. Student Table
CREATE TABLE Student (
    st_id INT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    gender CHAR(1) NOT NULL CHECK (gender IN ('M', 'F', 'O')),
    contact VARCHAR(15),
    dept VARCHAR(10),
    blood_group VARCHAR(5),
    dob DATE
);

-- 2. Doctor Table
CREATE TABLE Doctor (
    doc_id INT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    specialization VARCHAR(50) NOT NULL,
    contact VARCHAR(15) NOT NULL
);

-- 3. Token Table
CREATE TABLE Token (
    token_number INT NOT NULL,
    visit_date DATE NOT NULL DEFAULT (CURDATE()),
    visit_time TIME NOT NULL DEFAULT (CURTIME()),
    st_id INT NOT NULL,
    PRIMARY KEY (token_number, visit_date),
    CONSTRAINT fk_token_student
        FOREIGN KEY (st_id)
        REFERENCES Student(st_id)
        ON DELETE CASCADE
);

-- 4. Prescription Table
CREATE TABLE Prescription (
    p_id INT AUTO_INCREMENT PRIMARY KEY,
    p_date DATE NOT NULL DEFAULT (CURDATE()),
    st_id INT NOT NULL,
    doc_id INT NOT NULL,
    CONSTRAINT fk_presc_student
        FOREIGN KEY (st_id)
        REFERENCES Student(st_id)
        ON DELETE CASCADE,
    CONSTRAINT fk_presc_doctor
        FOREIGN KEY (doc_id)
        REFERENCES Doctor(doc_id)
        ON DELETE RESTRICT
);

-- 5. Pre_symptoms Table
CREATE TABLE Pre_symptoms (
    p_id INT NOT NULL,
    symptoms VARCHAR(100) NOT NULL,
    PRIMARY KEY (p_id, symptoms),
    CONSTRAINT fk_symptoms_presc
        FOREIGN KEY (p_id)
        REFERENCES Prescription(p_id)
        ON DELETE CASCADE
);

-- 6. Medicine Table
CREATE TABLE Medicine (
    medicine_name VARCHAR(100) NOT NULL,
    medicine_type VARCHAR(50) NOT NULL,
    manufacturer VARCHAR(100),
    PRIMARY KEY (medicine_name, medicine_type)
);

-- 7. Contain Table
CREATE TABLE Contain (
    p_id INT NOT NULL,
    medicine_name VARCHAR(100) NOT NULL,
    medicine_type VARCHAR(50) NOT NULL,
    dosage VARCHAR(50),
    PRIMARY KEY (p_id, medicine_name, medicine_type),
    CONSTRAINT fk_contain_presc
        FOREIGN KEY (p_id)
        REFERENCES Prescription(p_id)
        ON DELETE CASCADE,
    CONSTRAINT fk_contain_med
        FOREIGN KEY (medicine_name, medicine_type)
        REFERENCES Medicine(medicine_name, medicine_type)
        ON UPDATE CASCADE
);

-- 8. Test Table
CREATE TABLE Test (
    test_name VARCHAR(100) PRIMARY KEY,
    test_fee DECIMAL(10,2) NOT NULL DEFAULT 0.00
);

-- 9. Requires Table
CREATE TABLE Requires (
    p_id INT NOT NULL,
    test_name VARCHAR(100) NOT NULL,
    test_date DATE,
    test_result VARCHAR(255),
    PRIMARY KEY (p_id, test_name),
    CONSTRAINT fk_req_presc
        FOREIGN KEY (p_id)
        REFERENCES Prescription(p_id)
        ON DELETE CASCADE,
    CONSTRAINT fk_req_test
        FOREIGN KEY (test_name)
        REFERENCES Test(test_name)
        ON UPDATE CASCADE
);
 