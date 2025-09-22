CREATE TABLE attachments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ticketId UUID REFERENCES tickets(id) ON DELETE CASCADE,
    commentId UUID REFERENCES ticketComments(id) ON DELETE CASCADE,
    fileName VARCHAR(255) NOT NULL,
    fileUrl VARCHAR(500) NOT NULL, -- URL to the file (handled by separate file service)
    fileType VARCHAR(100), -- Type of file (document, image, video, etc.)
    uploadedBy VARCHAR(255) NOT NULL, -- External user ID from auth service
    metadata JSONB DEFAULT '{}', -- Additional flexible data
    createdAt TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT chk_attachment_reference CHECK (
        (ticketId IS NOT NULL AND commentId IS NULL) OR 
        (ticketId IS NULL AND commentId IS NOT NULL)
    )
);

CREATE INDEX idx_attachments_ticket_id ON attachments(ticketId);
CREATE INDEX idx_attachments_comment_id ON attachments(commentId);
CREATE INDEX idx_attachments_uploaded_by ON attachments(uploadedBy);
CREATE INDEX idx_attachments_metadata ON attachments USING GIN(metadata);