CREATE TABLE ticketComments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ticketId UUID NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    userId VARCHAR(255) NOT NULL, -- External user ID from auth service
    content TEXT NOT NULL,
    isInternal BOOLEAN DEFAULT false, -- internal notes vs customer-visible comments
    metadata JSONB DEFAULT '{}', -- Additional flexible data
    createdAt TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updatedAt TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_ticket_comments_ticket_id ON ticketComments(ticketId);
CREATE INDEX idx_ticket_comments_user_id ON ticketComments(userId);
CREATE INDEX idx_ticket_comments_created_at ON ticketComments(createdAt);
CREATE INDEX idx_ticket_comments_metadata ON ticketComments USING GIN(metadata);