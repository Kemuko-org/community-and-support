CREATE TABLE tickets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ticketNumber VARCHAR(20) UNIQUE NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    status VARCHAR(50) DEFAULT 'open' CHECK (status IN ('open', 'inProgress', 'waitingForCustomer', 'resolved', 'closed')),
    priority VARCHAR(50) DEFAULT 'medium' CHECK (priority IN ('low', 'medium', 'high', 'urgent')),
    type VARCHAR(50) DEFAULT 'general' CHECK (type IN ('general', 'technical', 'course', 'assignment', 'grading', 'platform', 'content')),
    studentId VARCHAR(255) NOT NULL, -- Student/User ID from auth service
    instructorId VARCHAR(255), -- Assigned instructor/support agent ID
    courseId VARCHAR(255), -- Course ID if ticket is course-related
    categoryId UUID REFERENCES categories(id) ON DELETE SET NULL,
    metadata JSONB DEFAULT '{}', -- Additional flexible data (course info, assignment details, etc.)
    createdAt TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updatedAt TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    resolvedAt TIMESTAMP WITH TIME ZONE,
    closedAt TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_tickets_status ON tickets(status);
CREATE INDEX idx_tickets_priority ON tickets(priority);
CREATE INDEX idx_tickets_type ON tickets(type);
CREATE INDEX idx_tickets_student_id ON tickets(studentId);
CREATE INDEX idx_tickets_instructor_id ON tickets(instructorId);
CREATE INDEX idx_tickets_course_id ON tickets(courseId);
CREATE INDEX idx_tickets_category_id ON tickets(categoryId);
CREATE INDEX idx_tickets_created_at ON tickets(createdAt);
CREATE INDEX idx_tickets_ticket_number ON tickets(ticketNumber);
CREATE INDEX idx_tickets_metadata ON tickets USING GIN(metadata);