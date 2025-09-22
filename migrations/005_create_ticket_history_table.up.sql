CREATE TABLE ticketHistory (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ticketId UUID NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    userId VARCHAR(255) NOT NULL, -- External user ID from auth service
    action VARCHAR(100) NOT NULL, -- 'created', 'statusChanged', 'assigned', 'priorityChanged', etc.
    oldValue TEXT,
    newValue TEXT,
    description TEXT,
    metadata JSONB DEFAULT '{}', -- Additional flexible data
    createdAt TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_ticket_history_ticket_id ON ticketHistory(ticketId);
CREATE INDEX idx_ticket_history_user_id ON ticketHistory(userId);
CREATE INDEX idx_ticket_history_created_at ON ticketHistory(createdAt);
CREATE INDEX idx_ticket_history_action ON ticketHistory(action);
CREATE INDEX idx_ticket_history_metadata ON ticketHistory USING GIN(metadata);