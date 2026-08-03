-- 1. Student Table
CREATE TABLE Student (
    st_id INT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    gender CHAR(1) NOT NULL CHECK (gender IN ('M', 'F', 'O')),
    contact VARCHAR(15),
    dept VARCHAR(10),
    dob DATE,
    blood_group VARCHAR(5)
    
);

-- 2. Token Table
CREATE TABLE Token (
    token_id INT NOT NULL,
    visit_date DATE NOT NULL,
    visit_time TIME NOT NULL,
    st_id INT NOT NULL,
    

    
    PRIMARY KEY (token_id, visit_date),
    CONSTRAINT fk_token_student 
        FOREIGN KEY (st_id) 
        REFERENCES Student(st_id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

-- 1. Doctor Table (Must be created before Prescription)
CREATE TABLE Doctor (
    doc_id INT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    specialization VARCHAR(50) NOT NULL,
    contact VARCHAR(15) NOT NULL
    
);

-- 2. Prescription Table
CREATE TABLE Prescription (
    p_id INT AUTO_INCREMENT PRIMARY KEY,
    p_date DATE DEFAULT (CURRENT_DATE),
    st_id INT NOT NULL,
    doc_id INT NOT NULL,
    
    CONSTRAINT fk_prescription_student 
        FOREIGN KEY (st_id) 
        REFERENCES Student(st_id)
        ON UPDATE CASCADE 
        ON DELETE CASCADE,
        
    CONSTRAINT fk_prescription_doctor 
        FOREIGN KEY (doc_id) 
        REFERENCES Doctor(doc_id)
        ON UPDATE CASCADE 
        ON DELETE RESTRICT
);

-- 3. Pre_symptoms Table (Multivalued attribute for Prescription)
CREATE TABLE Pre_symptoms (
    p_id INT NOT NULL,
    symptom VARCHAR(50) NOT NULL,
    
    PRIMARY KEY (p_id, symptom),
    CONSTRAINT fk_symptoms_prescription 
        FOREIGN KEY (p_id) 
        REFERENCES Prescription(p_id)
        ON DELETE CASCADE
);


-- 1. Master Medicine Catalog Table (#7)
CREATE TABLE Medicine (
    medicine_name VARCHAR(100) NOT NULL,
    medicine_type VARCHAR(50) NOT NULL,  -- e.g., 'Tablet', 'Syrup', 'Injection'
    manufacturer VARCHAR(100),
    
    PRIMARY KEY (medicine_name, medicine_type)
);

-- 2. Master Diagnostic Test Catalog Table (#8)
CREATE TABLE Test (
    test_name VARCHAR(100) PRIMARY KEY,
    test_fee DECIMAL(8,2) NOT NULL DEFAULT 0.00
);

-- 3. Prescription - Medicine Junction Table (#5 'contain')
CREATE TABLE Contain (
    p_id INT NOT NULL,
    medicine_name VARCHAR(100) NOT NULL,
    medicine_type VARCHAR(50) NOT NULL,
    dosage VARCHAR(50), -- Added field (e.g., '1-0-1 after food')
    
    PRIMARY KEY (p_id, medicine_name, medicine_type),
    
    CONSTRAINT fk_contain_prescription 
        FOREIGN KEY (p_id) 
        REFERENCES Prescription(p_id) 
        ON DELETE CASCADE,
        
    CONSTRAINT fk_contain_medicine 
        FOREIGN KEY (medicine_name, medicine_type) 
        REFERENCES Medicine(medicine_name, medicine_type) 
        ON UPDATE CASCADE
);

-- 4. Prescription - Test Junction Table (#6 'requires')
CREATE TABLE Requires (
    p_id INT NOT NULL,
    test_name VARCHAR(100) NOT NULL,
    test_date DATE,
    test_result VARCHAR(255),
    
    PRIMARY KEY (p_id, test_name),
    
    CONSTRAINT fk_requires_prescription 
        FOREIGN KEY (p_id) 
        REFERENCES Prescription(p_id) 
        ON DELETE CASCADE,
        
    CONSTRAINT fk_requires_test 
        FOREIGN KEY (test_name) 
        REFERENCES Test(test_name) 
        ON UPDATE CASCADE
);
